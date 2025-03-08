package graph

import (
	log "github.com/sirupsen/logrus"
	rabbitmqqueue "message-sender/internal/rabbitmq"
	"message-sender/microservices/subscriber/pkg/usecases"
)

type Subscriber struct {
	deps       SubscriberDependencies
	subscriber usecases.WSSubscriber
}

type SubscriberDependencies struct {
	logger *log.Logger
	queue  rabbitmqqueue.Queue
}

type SubscriberOptions struct {
	Logger *log.Logger
	Queue  rabbitmqqueue.Queue
}

func Install(opts SubscriberOptions) (*Subscriber, error) {
	subscriber := &Subscriber{
		deps: SubscriberDependencies{
			logger: opts.Logger,
			queue:  opts.Queue,
		},
	}
	return subscriber.install()
}

func (r *Subscriber) install() (*Subscriber, error) {
	r.deps.logger.SetLevel(log.InfoLevel)

	// provide use cases
	var newWSSubscriberOpts usecases.DefaultWSSubscriberOptions
	newWSSubscriberOpts.Logger = r.deps.logger
	newWSSubscriberOpts.URL = "ws://localhost:8081/echo"
	subscriber, err := usecases.NewWSSubscriber(newWSSubscriberOpts)
	if err != nil {
		return nil, err
	}

	r.subscriber = subscriber

	return r, nil
}

func (r *Subscriber) StartSubscriber() error {
	return r.subscriber.ReadMessage()
}
