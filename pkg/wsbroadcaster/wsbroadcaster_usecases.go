package wsbroadcaster

import (
	"errors"
	"fmt"
	"gig-assessment/internal/rabbitmq-queue"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type WSBroadcaster struct {
	queue   *rabbitmqqueue.RabbitMQQueue
	clients map[*websocket.Conn]bool
	mu      sync.Mutex
}

type NewWSBroadcasterOptions struct {
	Queue *rabbitmqqueue.RabbitMQQueue
}

func NewWSBroadcaster(opts NewWSBroadcasterOptions) (*WSBroadcaster, error) {
	if opts.Queue == nil {
		return nil, errors.New("option 'RabbitMQQueue' is mandatory")
	}

	return &WSBroadcaster{
		queue:   opts.Queue,
		clients: make(map[*websocket.Conn]bool),
		//ws:      nil,
	}, nil
}

func (wsb *WSBroadcaster) ReadQueueMessage() ([]byte, error) {
	for {
		message, err := wsb.queue.ReadMessage()
		if err != nil {
			msg := fmt.Sprintf("[WSBroadcaster] error reading message from the queue", err.Error())
			log.Println(msg)
			return nil, errors.New(msg)
		}

		go wsb.SendMessageToWS(message)
	}
}

func (wsb *WSBroadcaster) SendMessageToWS(message []byte) {
	for client := range wsb.clients {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("[WSBroadcaster] error sending message to subscribers:", err.Error())
		}
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
	log.Println("[WSBroadcaster] New client connected")
}

func (wsb *WSBroadcaster) DisconnectClient(conn *websocket.Conn) {
	wsb.mu.Lock()
	delete(wsb.clients, conn)
	wsb.mu.Unlock()
	fmt.Println("[WSBroadcaster] Client disconnected")
}
