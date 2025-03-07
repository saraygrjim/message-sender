package microservices

import (
	"errors"
	"fmt"
	"gig-assessment/internal/rabbitmq-queue"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"sync"
)

type WSBroadcaster struct {
	queue   *rabbitmqqueue.RabbitMQQueue
	clients map[*websocket.Conn]bool
	mu      sync.Mutex
	logger  *log.Logger
}

type NewWSBroadcasterOptions struct {
	Queue  *rabbitmqqueue.RabbitMQQueue
	Logger *log.Logger
}

func NewWSBroadcaster(opts NewWSBroadcasterOptions) (*WSBroadcaster, error) {
	if opts.Queue == nil {
		return nil, errors.New("option 'RabbitMQQueue' is mandatory")
	}
	if opts.Logger == nil {
		return nil, errors.New("option 'RabbitMQLogger' is mandatory")
	}

	return &WSBroadcaster{
		queue:   opts.Queue,
		clients: make(map[*websocket.Conn]bool),
		logger:  opts.Logger,
	}, nil
}

func (wsb *WSBroadcaster) ReadQueueMessage() ([]byte, error) {
	for {
		wsb.logger.Info("[WSBroadcaster] reading queue messages")

		message, err := wsb.queue.ReadMessage()
		if err != nil {
			wsb.logger.WithFields(log.Fields{
				ErrorTag: err.Error(),
			}).Fatal("[WSBroadcaster] error reading queue messages")

			return nil, errors.New(fmt.Sprintf("error reading queue messages: %s", err.Error()))
		}

		go wsb.SendMessageToWS(message)
	}
}

func (wsb *WSBroadcaster) SendMessageToWS(message []byte) {
	wsb.logger.Info("[WSBroadcaster] sending messages to subscribers")

	var hasErrors bool

	for client := range wsb.clients {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			hasErrors = true
			wsb.logger.WithFields(log.Fields{
				SubscriberTag: client.RemoteAddr(),
				MessageTag:    message,
				ErrorTag:      err.Error(),
			}).Error("[WSBroadcaster]  error sending message to subscriber")
		}
	}

	if !hasErrors {
		wsb.logger.Debug("[WSBroadcaster] message succesfully sent to all the subscribers")
	}

	return
}

func (wsb *WSBroadcaster) ReadClientMessage(conn *websocket.Conn) {
	wsb.NewClient(conn)
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			wsb.DisconnectClient(conn)
			break
		}
	}
}

func (wsb *WSBroadcaster) NewClient(conn *websocket.Conn) {
	wsb.mu.Lock()
	wsb.clients[conn] = true
	wsb.mu.Unlock()

	wsb.logger.WithFields(log.Fields{
		SubscriberTag: conn.RemoteAddr(),
	}).Info("[WSBroadcaster]  new client connected")

}

func (wsb *WSBroadcaster) DisconnectClient(conn *websocket.Conn) {
	wsb.mu.Lock()
	delete(wsb.clients, conn)
	wsb.mu.Unlock()

	wsb.logger.WithFields(log.Fields{
		SubscriberTag: conn.RemoteAddr(),
	}).Info("[WSBroadcaster] client disconnected")
}
