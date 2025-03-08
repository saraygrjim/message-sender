package http

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"message-sender/internal/websocket"
	"message-sender/microservices/broadcaster/pkg/usecases"
	"net/http"
)

type WSBroadcasterAdapterOptions struct {
	Logger        *log.Logger
	WSBroadcaster usecases.WSBroadcaster
}

func ProvideWSBroadcasterAdapter(opts WSBroadcasterAdapterOptions) (*WSBroadcasterAdapter, error) {
	if opts.Logger == nil {
		return nil, errors.New("option 'Logger' is mandatory")
	}
	if opts.WSBroadcaster == nil {
		return nil, errors.New("option 'WSBroadcaster' is mandatory")
	}

	return &WSBroadcasterAdapter{
		logger:        opts.Logger,
		wsBroadcaster: opts.WSBroadcaster,
	}, nil
}

var _ usecases.WSBroadcasterAdapter = (*WSBroadcasterAdapter)(nil)

type WSBroadcasterAdapter struct {
	logger        *log.Logger
	wsBroadcaster usecases.WSBroadcaster
}

func (a WSBroadcasterAdapter) ServeWebsocketBroadcasterHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.NewConnection(w, r, a.logger)
	if err != nil {
		panic(err)
	}

	go a.wsBroadcaster.ReadClientMessage(conn)
}
