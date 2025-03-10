// Microservice 2
package usecases

import (
	"github.com/go-errors/errors"
	log "github.com/sirupsen/logrus"
	rabbitmqqueue "message-sender/internal/rabbitmq"
	"message-sender/internal/websocket"
	"sync"
)

type WSBroadcaster interface {
	ReadQueueMessage() ([]byte, error)
	SendMessageToWS(message []byte) error
	ReadClientMessage(conn *websocket.WSConnection) error
	NewClient(conn *websocket.WSConnection)
	DisconnectClient(conn *websocket.WSConnection)
}
type DefaultWSBroadcasterOptions struct {
	Queue  rabbitmqqueue.Queue
	Logger *log.Logger
}

func NewWSBroadcaster(opts DefaultWSBroadcasterOptions) (*DefaultWSBroadcaster, error) {
	if opts.Queue == nil {
		return nil, errors.New("option 'Queue' is mandatory")
	}
	if opts.Logger == nil {
		return nil, errors.New("option 'Logger' is mandatory")
	}

	return &DefaultWSBroadcaster{
		queue:   opts.Queue,
		clients: make(map[*websocket.WSConnection]bool),
		logger:  opts.Logger,
	}, nil
}

var _ WSBroadcaster = (*DefaultWSBroadcaster)(nil)

type DefaultWSBroadcaster struct {
	queue   rabbitmqqueue.Queue
	clients map[*websocket.WSConnection]bool
	mu      sync.Mutex
	logger  *log.Logger
}

func (wsb *DefaultWSBroadcaster) ReadQueueMessage() ([]byte, error) {
	for {
		wsb.logger.Info("[WSBroadcaster] reading queue messages")

		message, err := wsb.queue.ReadMessage()
		if err != nil {
			wsb.logger.WithFields(log.Fields{
				ErrorTag: err.Error(),
			}).Fatal("[WSBroadcaster] error reading queue messages")

			return nil, ErrReadingQueueMessages
		}

		go wsb.SendMessageToWS(message)
	}
}

func (wsb *DefaultWSBroadcaster) SendMessageToWS(message []byte) error {
	wsb.logger.Info("[WSBroadcaster] sending messages to subscribers")

	var hasErrors bool

	for client := range wsb.clients {
		err := client.WriteMessage(message)
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
		wsb.logger.Debug("[WSBroadcaster] message successfully sent to all the subscribers")
		return nil
	} else {
		return ErrSendingMessageToWebsocket
	}
}

func (wsb *DefaultWSBroadcaster) ReadClientMessage(conn *websocket.WSConnection) error {
	wsb.NewClient(conn)
	for {
		_, err := conn.ReadMessage()
		if err != nil {
			wsb.DisconnectClient(conn)
			return ErrClientDisconnected
		}
	}
}

func (wsb *DefaultWSBroadcaster) NewClient(conn *websocket.WSConnection) {
	wsb.mu.Lock()
	wsb.clients[conn] = true
	wsb.mu.Unlock()

	wsb.logger.WithFields(log.Fields{
		SubscriberTag: conn.RemoteAddr(),
	}).Info("[WSBroadcaster]  new client connected")

}

func (wsb *DefaultWSBroadcaster) DisconnectClient(conn *websocket.WSConnection) {
	wsb.mu.Lock()
	delete(wsb.clients, conn)
	wsb.mu.Unlock()

	wsb.logger.WithFields(log.Fields{
		SubscriberTag: conn.RemoteAddr(),
	}).Info("[WSBroadcaster] client disconnected")
}

// log tags
const (
	MessageTag    = "message"
	ErrorTag      = "error"
	SubscriberTag = "subscriber"
)

// errors

var ErrReadingQueueMessages = errors.Errorf("error reading queue messages")
var ErrSendingMessageToWebsocket = errors.Errorf("error sending message to websocket")
var ErrClientDisconnected = errors.Errorf("client disconnected")
