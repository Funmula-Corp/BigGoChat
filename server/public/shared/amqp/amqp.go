package amqp

import (
	"context"
	"errors"
	"net"
	"runtime"
	"time"

	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/mlog"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	amqpTimeOut     = 5 * time.Second
	maxRetryWaiting = 30 * time.Second
)

var (
	numWorker = runtime.NumCPU() * 8
)

type Cmd int

const (
	commandReserved Cmd = iota // don't use 0
	commandReconnect
)

type amqpValue string

const (
	amqpExchanges amqpValue = "exchanges"
)

// Connect to server with const timeout, run forever until canceled
func Connect(ctx context.Context, url string) (*amqp.Connection, chan *amqp.Error, error) {
	currentRetryWaiting := time.Second
	config := amqp.Config{
		Locale: "en_US",
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, amqpTimeOut)
		},
	}

	for {
		select {
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		default:
			client, err := amqp.DialConfig(url, config)
			if err != nil {
				mlog.Error("AMQPClient: cannot connect amqp server", mlog.Err(err))
				time.Sleep(min(currentRetryWaiting, maxRetryWaiting))
				currentRetryWaiting *= 2
				continue
			}
			mlog.Info("AMQPClient: remote server connected")
			return client, client.NotifyClose(make(chan *amqp.Error)), nil
		}
	}
}

func getChannel(ctx context.Context, client *amqp.Connection, onChannelCreated func(*amqp.Channel) error) (*amqp.Channel, chan *amqp.Error, error) {
	currentRetryWaiting := time.Second
	for {
		select {
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		default:
			channel, err := client.Channel()
			if err != nil {
				if errors.Is(err, amqp.ErrClosed) {
					// connection closed
					return nil, nil, err
				}
				mlog.Info("AMQPClient: cannot create channel", mlog.Err(err))
				time.Sleep(min(currentRetryWaiting, maxRetryWaiting))
				currentRetryWaiting *= 2
				continue
			}
			// mlog.Info("AMQPClient: channel created")
			if onChannelCreated != nil {
				err = onChannelCreated(channel)
			}
			return channel, channel.NotifyClose(make(chan *amqp.Error)), err
		}
	}
}

type AMQPMessage struct {
	Exchange string
	Key      string
	Headers  map[string]interface{}
	Body     []byte
}

type AMQPClient struct {
	url       string
	client    *amqp.Connection
	ctx       context.Context // root context
	cancel    func()          // root cancel
	queue     chan AMQPMessage
	cmd       chan Cmd
	exchanges []Exchange
}

type Exchange struct {
	Name string
	Type string
}

// **MUST** call this to create AMQPClient service
func MakeAMQPClient(url string, exchanges []Exchange) *AMQPClient {
	ctx, cancel := context.WithCancel(context.Background())
	amqpClient := &AMQPClient{
		url: url,
		// buffer already in PushNotificationsHub
		queue:     make(chan AMQPMessage),
		ctx:       ctx,
		cancel:    cancel,
		cmd:       make(chan Cmd),
		exchanges: exchanges,
	}

	go amqpClient.supervisor()

	return amqpClient
}

func (a *AMQPClient) Publish(message AMQPMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), amqpTimeOut)
	defer cancel()
	select {
	case a.queue <- message:
		return nil
	case <-a.ctx.Done():
		return a.ctx.Err()
	case <-ctx.Done():
		// in case of time out
		return ctx.Err()
	}
}

// **unsafe** expose channel for other usage
func (a *AMQPClient) Channel() (*amqp.Channel, error) {
	return a.client.Channel()
}

func (a *AMQPClient) Shutdown() {
	a.cancel()
}

// connect to new url without drop data
func (a *AMQPClient) SwitchToNewServer(url string) error {
	a.url = url
	select {
	case a.cmd <- commandReconnect:
	case <-a.ctx.Done():
		return a.ctx.Err()
	}
	return nil
}

func (a *AMQPClient) supervisor() {
	// startup
	client, notifyClose, err := Connect(a.ctx, a.url)
	if err != nil {
		// only root cancel can cause err
		return
	}
	a.client = client

	createWorkers := func(ctx context.Context, c int) (func(), <-chan error) {
		sub, cancel := context.WithCancel(context.WithValue(ctx, amqpExchanges, a.exchanges))
		ch := make(chan error)
		for i := 0; i < c; i++ {
			go worker(sub, a.client, a.queue, ch)
		}
		return cancel, ch
	}

	cancelWorkers, workerErr := createWorkers(a.ctx, numWorker)

	// watch
	for {
		select {
		case notify := <-notifyClose:
			if notify == nil {
				mlog.Info("AMQPClient: connection closed")
			} else {
				mlog.Info("AMQPClient: connection closed", mlog.Err(notify))
			}
			cancelWorkers()
			a.client, notifyClose, err = Connect(a.ctx, a.url)
			if err != nil {
				// root canceled, go to ctx.Done()
				continue
			}
			cancelWorkers, workerErr = createWorkers(a.ctx, numWorker)
		case cmd := <-a.cmd:
			if cmd == commandReconnect {
				a.client.Close()
			}
		case <-workerErr:
			// ...
		case <-a.ctx.Done():
			// full shutdown
			if a.client != nil {
				a.client.Close()
			}
			return
		}
	}
}

func worker(ctx context.Context, client *amqp.Connection, msgCh chan AMQPMessage, errCh chan<- error) {
	handleErr := func(err error) (exit bool) {
		if err == nil {
			return false
		} else if errors.Is(err, amqp.ErrClosed) && client.IsClosed() {
			return true
		} else if errors.Is(err, context.Canceled) {
			return true
		}
		mlog.Error("AMQPClient: worker failed", mlog.Err(err))
		select {
		case errCh <- err:
		case <-ctx.Done():
		}
		return false
	}

	initChannel := func(channel *amqp.Channel) error {
		// enable confirm mode
		err := channel.Confirm(false)
		if err != nil {
			return err
		}

		// create exchange
		for _, exchange := range ctx.Value(amqpExchanges).([]Exchange) {
			if exchange.Type == "" {
				// default type
				err = channel.ExchangeDeclare(exchange.Name, amqp.ExchangeTopic, true, false, false, false, nil)
			} else {
				err = channel.ExchangeDeclare(exchange.Name, exchange.Type, true, false, false, false, nil)
			}

			if err != nil {
				return err
			}
		}
		return nil
	}

	channel, notify, err := getChannel(ctx, client, initChannel)
	if handleErr(err) {
		return
	}

	for {
		select {
		case message := <-msgCh:
			ok, err := publish(ctx, channel, message)
			if !ok || err != nil {
				go func() { msgCh <- message }()
			}
			handleErr(err)
		case <-notify:
			channel, notify, err = getChannel(ctx, client, initChannel)
			if handleErr(err) {
				return
			}
		case <-ctx.Done():
			channel.Close()
			return
		}
	}
}

func publish(ctx context.Context, channel *amqp.Channel, message AMQPMessage) (bool, error) {
	confirm, err := channel.PublishWithDeferredConfirm(message.Exchange, message.Key, false, false, amqp.Publishing{
		Headers: message.Headers,
		Body:    message.Body,
	})
	if err != nil {
		return false, err
	}

	sub, cancel := context.WithTimeout(ctx, amqpTimeOut)
	defer cancel()

	return confirm.WaitContext(sub)
}
