// Microservice 1
package microservices

import (
	"bytes"
	"errors"
	"fmt"
	"gig-assessment/internal/rabbitmq-queue"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type WSReceiver struct {
	queue  *rabbitmqqueue.RabbitMQQueue
	logger *log.Logger
}

type NewWSReceiverOptions struct {
	Queue  *rabbitmqqueue.RabbitMQQueue
	Logger *log.Logger
}

func NewWSReceiver(opts NewWSReceiverOptions) (*WSReceiver, error) {
	if opts.Queue == nil {
		return nil, errors.New("option 'RabbitMQQueue' is mandatory")
	}
	if opts.Logger == nil {
		return nil, errors.New("option 'RabbitMQLogger' is mandatory")
	}

	return &WSReceiver{
		queue:  opts.Queue,
		logger: opts.Logger,
	}, nil
}

func (wsr *WSReceiver) ReadMessage(conn *websocket.Conn) {
	wsr.logger.Info("[WSReceiver] reading message")
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				wsr.logger.WithFields(log.Fields{
					ErrorTag: err.Error(),
				}).Fatal("[WSReceiver] channel unexpectedly closed")
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

func (wsr *WSReceiver) SendMessageToQueue(message []byte) {
	fmt.Println("non format message " + string(message))                             //todo: delete
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
