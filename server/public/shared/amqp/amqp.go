package amqp

import (
	"context"
	"errors"
	"runtime"
	"sync"
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

type AMQPMessage struct {
	Exchange string
	Key      string
	Headers  map[string]interface{}
	Body     []byte
}

type AMQPClient struct {
	url    string
	client *amqp.Connection
	ctx    context.Context
	cancel func()
	queue  chan AMQPMessage
	lock   *sync.Mutex
}

func MakeAMQPClient(url string) *AMQPClient {
	amqpClient := &AMQPClient{
		url: url,
		// buffer already in PushNotificationsHub
		queue: make(chan AMQPMessage),
		lock:  new(sync.Mutex),
	}
	go amqpClient.reconnect()

	return amqpClient
}

func (a *AMQPClient) Publish(message AMQPMessage) {
	a.queue <- message
}

// expose channel for other usage
func (a *AMQPClient) Channel() (*amqp.Channel, error) {
	return a.client.Channel()
}

func (a *AMQPClient) Shutdown() {
	if a.cancel != nil {
		a.cancel()
	}
	if a.client != nil {
		a.client.Close()
	}
}

func (a *AMQPClient) reconnect() {
	if !a.lock.TryLock() {
		return
	}
	defer a.lock.Unlock()

	a.Shutdown()

	currentRetryWaiting := time.Second

	// retry connect until amqp connected
	for {
		client, err := amqp.Dial(a.url)
		if err != nil {
			mlog.Error("AMQPClient: cannot connect amqp server", mlog.Err(err))
			time.Sleep(min(currentRetryWaiting, maxRetryWaiting))
			currentRetryWaiting *= 2
			continue
		}
		a.client = client
		a.ctx, a.cancel = context.WithCancel(context.Background())

		// start workers
		for i := 0; i < numWorker; i++ {
			go a.worker()
		}

		break
	}

}

// can only call by worker
func (a *AMQPClient) handleErr(err error) {
	mlog.Error("AMQPClient: worker failed", mlog.Err(err))

	// should always reconnect or just amqp server down?
	if errors.Is(err, amqp.ErrClosed) {
		a.reconnect()
	}
}

func (a *AMQPClient) worker() {
	channel, err := a.client.Channel()
	if err != nil {
		defer a.handleErr(err)
		return
	}

	for {
		select {
		case message := <-a.queue:
			if err := publish(a.ctx, channel, message); err != nil {
				defer a.handleErr(err)
				// put it back
				go func() { a.queue <- message }()
				return
			}
		case <-a.ctx.Done():
			return
		}
	}
}

func publish(parent context.Context, channel *amqp.Channel, message AMQPMessage) error {
	ctx, cancel := context.WithTimeout(parent, amqpTimeOut)
	defer cancel()

	err := channel.PublishWithContext(ctx, message.Exchange, message.Key, false, false, amqp.Publishing{
		Headers: message.Headers,
		Body:    message.Body,
	})

	return err
}
