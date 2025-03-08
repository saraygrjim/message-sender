package http

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"message-sender/internal/websocket"
	"message-sender/microservices/receiver/pkg/usecases"
	"net/http"
)

type WSReceiverAdapterOptions struct {
	Logger     *log.Logger
	WSReceiver usecases.WSReceiver
}

func ProvideWSReceiverAdapter(opts WSReceiverAdapterOptions) (*WSReceiverAdapter, error) {
	if opts.Logger == nil {
		return nil, errors.New("option 'Logger' is mandatory")
	}
	if opts.WSReceiver == nil {
		return nil, errors.New("option 'WSReceiver' is mandatory")
	}

	return &WSReceiverAdapter{
		logger:     opts.Logger,
		wsReceiver: opts.WSReceiver,
	}, nil
}

var _ usecases.WSReceiverAdapter = (*WSReceiverAdapter)(nil)

type WSReceiverAdapter struct {
	logger     *log.Logger
	wsReceiver usecases.WSReceiver
}

func (a WSReceiverAdapter) ServeWebsocketReceiverHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.NewConnection(w, r, a.logger)
	if err != nil {
		panic(err)
	}

	go a.wsReceiver.ReadMessage(conn)
}
