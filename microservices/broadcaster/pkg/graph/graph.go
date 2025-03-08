package graph

import (
	log "github.com/sirupsen/logrus"
	rabbitmqqueue "message-sender/internal/rabbitmq"
	httpAdapter "message-sender/microservices/broadcaster/pkg/adapters/input/http"
	"message-sender/microservices/broadcaster/pkg/usecases"
	"net/http"
)

type Broadcaster struct {
	deps        BroadcasterDependencies
	broadcaster usecases.WSBroadcaster
	httpAdapter usecases.WSBroadcasterAdapter
}

type BroadcasterDependencies struct {
	logger *log.Logger
	queue  rabbitmqqueue.Queue
}

type BroadcasterOptions struct {
	Logger *log.Logger
	Queue  rabbitmqqueue.Queue
}

func Install(opts BroadcasterOptions) (*Broadcaster, error) {
	broadcaster := &Broadcaster{
		deps: BroadcasterDependencies{
			logger: opts.Logger,
			queue:  opts.Queue,
		},
	}
	return broadcaster.install()
}

func (r *Broadcaster) install() (*Broadcaster, error) {
	r.deps.logger.SetLevel(log.InfoLevel)

	// provide use cases
	var newWSBroadcasterOpts usecases.DefaultWSBroadcasterOptions
	newWSBroadcasterOpts.Queue = r.deps.queue
	newWSBroadcasterOpts.Logger = r.deps.logger
	broadcaster, err := usecases.NewWSBroadcaster(newWSBroadcasterOpts)
	if err != nil {
		return nil, err
	}

	r.broadcaster = broadcaster

	// provide adapters
	var wsBroadcasterAdapterOpts httpAdapter.WSBroadcasterAdapterOptions
	wsBroadcasterAdapterOpts.WSBroadcaster = r.broadcaster
	wsBroadcasterAdapterOpts.Logger = r.deps.logger
	adapter, err := httpAdapter.ProvideWSBroadcasterAdapter(wsBroadcasterAdapterOpts)
	if err != nil {
		return nil, err
	}
	r.httpAdapter = adapter

	// register handlers
	http.HandleFunc("/echo", r.httpAdapter.ServeWebsocketBroadcasterHTTP)

	return r, nil
}

func (r *Broadcaster) StartWebsocketBroadcasterServer() {
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		panic(err)
	}
}

func (r *Broadcaster) StartBroadcaster() {
	_, err := r.broadcaster.ReadQueueMessage()
	if err != nil {
		panic(err)
	}
}
