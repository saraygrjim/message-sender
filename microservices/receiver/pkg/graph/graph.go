package graph

import (
	log "github.com/sirupsen/logrus"
	rabbitmqqueue "message-sender/internal/rabbitmq"
	httpAdapter "message-sender/microservices/receiver/pkg/adapters/input/http"
	"message-sender/microservices/receiver/pkg/usecases"
	"net/http"
)

type Receiver struct {
	deps        ReceiverDependencies
	receiver    usecases.WSReceiver
	httpAdapter usecases.WSReceiverAdapter
}

type ReceiverDependencies struct {
	logger *log.Logger
	queue  rabbitmqqueue.Queue
}

type ReceiverOptions struct {
	Logger *log.Logger
	Queue  rabbitmqqueue.Queue
}

func Install(opts ReceiverOptions) (*Receiver, error) {
	receiver := &Receiver{
		deps: ReceiverDependencies{
			logger: opts.Logger,
			queue:  opts.Queue,
		},
	}
	return receiver.install()
}

func (r *Receiver) install() (*Receiver, error) {
	r.deps.logger.SetLevel(log.InfoLevel)

	// provide use cases
	var newWSReceiverOpts usecases.DefaultWSReceiverOptions
	newWSReceiverOpts.Queue = r.deps.queue
	newWSReceiverOpts.Logger = r.deps.logger
	receiver, err := usecases.NewWSReceiver(newWSReceiverOpts)
	if err != nil {
		return nil, err
	}

	r.receiver = receiver

	// provide adapters
	var wsReceiverAdapterOpts httpAdapter.WSReceiverAdapterOptions
	wsReceiverAdapterOpts.WSReceiver = r.receiver
	wsReceiverAdapterOpts.Logger = r.deps.logger
	adapter, err := httpAdapter.ProvideWSReceiverAdapter(wsReceiverAdapterOpts)
	if err != nil {
		return nil, err
	}
	r.httpAdapter = adapter

	// register handlers
	http.HandleFunc("/ws", r.httpAdapter.ServeWebsocketReceiverHTTP)

	return r, nil
}

func (r *Receiver) StartWebsocketReceiverServer() {
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
