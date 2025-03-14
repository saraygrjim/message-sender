package usecases

import (
	"github.com/go-errors/errors"
	log "github.com/sirupsen/logrus"
	"message-sender/internal/websocket"
)

type WSSubscriber interface {
	ReadMessage() error
}

type DefaultWSSubscriberOptions struct {
	Logger *log.Logger
	URL    string
}

func NewWSSubscriber(opts DefaultWSSubscriberOptions) (*DefaultWSSubscriber, error) {
	ws, err := websocket.NewClient(opts.URL, opts.Logger)
	if err != nil {
		return nil, err
	}

	return &DefaultWSSubscriber{
		logger:     opts.Logger,
		connection: ws,
	}, nil

}

var _ WSSubscriber = (*DefaultWSSubscriber)(nil)

type DefaultWSSubscriber struct {
	logger     *log.Logger
	connection *websocket.WSConnection
}

func (s DefaultWSSubscriber) ReadMessage() error {
	defer s.connection.Close()
	s.logger.Info("📡 Connected to websocket. Waiting for messages...")
	for {
		msg, err := s.connection.ReadMessage()
		if err != nil {
			s.logger.WithFields(log.Fields{
				ErrorTag: err.Error(),
			}).Error("Error reading message")
			e := ErrReadingWSMessage
			return e
		}
		s.logger.WithFields(log.Fields{MessageTag: string(msg)}).Info("📩 Message received")
	}
}

var ErrReadingWSMessage = errors.Errorf("error reading message from websocket")

const (
	MessageTag = "message"
	ErrorTag   = "error"
)
