// Microservice 1
package usecases

import (
	"bytes"
	"errors"
	log "github.com/sirupsen/logrus"
	"message-sender/internal/rabbitmq"
	"message-sender/internal/websocket"
)

type WSReceiver interface {
	ReadMessage(conn *websocket.Connection)
	SendMessageToQueue(message []byte)
}
type DefaultWSReceiverOptions struct {
	Queue  rabbitmqqueue.Queue
	Logger *log.Logger
}

var _ WSReceiver = (*DefaultWSReceiver)(nil)

type DefaultWSReceiver struct {
	queue  rabbitmqqueue.Queue
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

func (wsr *DefaultWSReceiver) ReadMessage(conn *websocket.Connection) {
	wsr.logger.Info("[WSReceiver] channel open. Reading messages...")
	defer conn.Close()
	for {
		message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				wsr.logger.WithFields(log.Fields{
					ErrorTag: err.Error(),
				}).Info("[WSReceiver] channel closed")
				return
			}

			wsr.logger.WithFields(log.Fields{
				ErrorTag: err.Error(),
			}).Error("[WSReceiver] there was an error reading the message")
		}
		if message == nil {
			continue
		}
		wsr.logger.Debug("[WSReceiver] message read")
		go wsr.SendMessageToQueue(message)

	}
}

func (wsr *DefaultWSReceiver) SendMessageToQueue(message []byte) {
	message = bytes.TrimSpace(bytes.Replace(message, []byte("\n"), []byte(" "), -1)) //todo: I am assuming that this is parsing the message, improve it with a function

	wsr.logger.WithFields(log.Fields{
		MessageTag: message,
	}).Info("[WSReceiver] sending message to queue")

	err := wsr.queue.SendMessage(message)
	if err != nil {
		wsr.logger.WithFields(log.Fields{
			MessageTag: message,
			ErrorTag:   err.Error(),
		}).Error("[WSReceiver] there was an error writing message '%s' into queue")
		return
	}

	wsr.logger.WithFields(log.Fields{
		MessageTag: message,
	}).Debug("[WSReceiver] message sent to queue")
}
