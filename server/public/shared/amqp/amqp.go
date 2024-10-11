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

	ErrSwitchNewServer = errors.New("amqp switch to new server")
)

// connect to server with const timeout, run forever until canceled
func connect(ctx context.Context, url string) (*amqp.Connection, chan *amqp.Error, error) {
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

type AMQPMessage struct {
	Exchange string
	Key      string
	Headers  map[string]interface{}
	Body     []byte
}

type AMQPClient struct {
	url    string
	client *amqp.Connection
	ctx    context.Context // root context
	cancel func()          // root cancel
	queue  chan AMQPMessage
	err    chan error
}

// **MUST** call this to create AMQPClient service
func MakeAMQPClient(url string) *AMQPClient {
	ctx, cancel := context.WithCancel(context.Background())
	amqpClient := &AMQPClient{
		url: url,
		// buffer already in PushNotificationsHub
		queue:  make(chan AMQPMessage),
		ctx:    ctx,
		cancel: cancel,
		err:    make(chan error),
	}

	go amqpClient.supervisor()

	return amqpClient
}

func (a *AMQPClient) Publish(message AMQPMessage) error {
	select {
	case a.queue <- message:
		return nil
	case <-a.ctx.Done():
		return a.ctx.Err()
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
	case a.err <- ErrSwitchNewServer:
	case <-a.ctx.Done():
		return a.ctx.Err()
	}
	return nil
}

func (a *AMQPClient) supervisor() {
	// startup
	client, notifyClose, err := connect(a.ctx, a.url)
	if err != nil {
		// only root cancel can cause err
		return
	}
	a.client = client

	createWorkers := func(ctx context.Context, c int) func() {
		sub, cancel := context.WithCancel(ctx)
		for i := 0; i < c; i++ {
			go a.worker(sub)
		}
		return cancel
	}

	cancelWorkers := createWorkers(a.ctx, numWorker)

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
			a.client, notifyClose, err = connect(a.ctx, a.url)
			if err != nil {
				// root canceled
				return
			}
			cancelWorkers = createWorkers(a.ctx, numWorker)
		case err := <-a.err:
			if errors.Is(err, ErrSwitchNewServer) {
				a.client.Close()
			}
		case <-a.ctx.Done():
			// full shutdown
			if a.client != nil {
				a.client.Close()
			}
		}
	}
}

func (a *AMQPClient) worker(ctx context.Context) {
	handleErr := func(err error) {
		mlog.Error("AMQPClient: worker failed", mlog.Err(err))
		select {
		case a.err <- err:
		case <-ctx.Done():
		}
	}

	channel, err := a.client.Channel()
	if err != nil {
		handleErr(err)
		return
	}
	defer channel.Close()

	// enable confirm mode
	if err := channel.Confirm(false); err != nil {
		handleErr(err)
	}

	for {
		select {
		case message := <-a.queue:
			ok, err := publish(ctx, channel, message)
			if !ok || err != nil {
				go func() { a.queue <- message }()
			}
			if err != nil {
				handleErr(err)
			}
		case <-ctx.Done():
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
