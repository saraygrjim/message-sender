// Microservice 1
package usecases

import (
	//"errors"
	"github.com/go-errors/errors"
	log "github.com/sirupsen/logrus"
	"message-sender/internal/rabbitmq"
	"message-sender/internal/websocket"
)

type WSReceiver interface {
	ReadMessage(conn websocket.Connection) error
	SendMessageToQueue(message []byte) error
}
type DefaultWSReceiverOptions struct {
	Queue  rabbitmq.Queue
	Logger *log.Logger
}

var _ WSReceiver = (*DefaultWSReceiver)(nil)

type DefaultWSReceiver struct {
	queue  rabbitmq.Queue
	logger *log.Logger
}

func NewWSReceiver(opts DefaultWSReceiverOptions) (*DefaultWSReceiver, error) {
	if opts.Queue == nil {
		return nil, errors.New("option 'Queue' is mandatory")
	}
	if opts.Logger == nil {
		return nil, errors.New("option 'Logger' is mandatory")
	}

	return &DefaultWSReceiver{
		queue:  opts.Queue,
		logger: opts.Logger,
	}, nil
}

func (wsr *DefaultWSReceiver) ReadMessage(conn websocket.Connection) error {
	wsr.logger.Info("[WSReceiver] channel open. Reading messages...")
	defer conn.Close()
	for {
		message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				wsr.logger.WithFields(log.Fields{
					ErrorTag: err.Error(),
				}).Info("[WSReceiver] channel closed")
				return nil
			}

			wsr.logger.WithFields(log.Fields{
				ErrorTag: err.Error(),
			}).Error("[WSReceiver] there was an error reading the message")

			return ErrReadingWSMessage
		}
		if message == nil {
			continue
		}
		wsr.logger.Debug("[WSReceiver] message read")
		err = wsr.SendMessageToQueue(message)
		if err != nil {
			return err
		}

	}
}

func (wsr *DefaultWSReceiver) SendMessageToQueue(message []byte) error {
	wsr.logger.WithFields(log.Fields{
		MessageTag: string(message),
	}).Info("[WSReceiver] sending message to queue")

	err := wsr.queue.SendMessage(message)
	if err != nil {
		wsr.logger.WithFields(log.Fields{
			MessageTag: string(message),
			ErrorTag:   err.Error(),
		}).Error("[WSReceiver] there was an error writing message into queue")
		return ErrWritingMessageInQueue
	}

	wsr.logger.WithFields(log.Fields{
		MessageTag: string(message),
	}).Debug("[WSReceiver] message sent to queue")

	return nil
}

// log tags
const (
	MessageTag    = "message"
	ErrorTag      = "error"
	SubscriberTag = "subscriber"
)

// errors
var ErrReadingWSMessage = errors.Errorf("error reading message from websocket")
var ErrWritingMessageInQueue = errors.Errorf("error writing message to queue")
