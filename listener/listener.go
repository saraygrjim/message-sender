package listener

import (
	"errors"
	"fmt"
	"gig-assessment/rabbitmq-queue"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

type Listener struct {
	queue *rabbitmqqueue.Queue
	//ws      *websocket.Conn
	clients map[*websocket.Conn]bool
	mu      sync.Mutex
}

type NewListenerOptions struct {
	Queue *rabbitmqqueue.Queue
}

func NewListener(opts NewListenerOptions) (*Listener, error) {
	if opts.Queue == nil {
		return nil, errors.New("option 'Queue' is mandatory")
	}

	return &Listener{
		queue:   opts.Queue,
		clients: make(map[*websocket.Conn]bool),
		//ws:      nil,
	}, nil
}

func (p *Listener) ReadMessage() ([]byte, error) {
	for {
		message, err := p.queue.ReadMessage()
		if err != nil {
			msg := fmt.Sprintf("[Listener] error reading message from the queue", err.Error())
			log.Println(msg)
			return nil, errors.New(msg)
		}

		err = p.SendMessageToWS(message)
		if err != nil {
			return nil, err
		}
	}
}

func (p *Listener) SendMessageToWS(message []byte) error {
	for client := range p.clients {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			msg := fmt.Sprintf("[Listener] error sending message to subscribers", err.Error())
			log.Println("write:", err)
			return errors.New(msg)
		}
	}

	return nil
}

func (p *Listener) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("[Listener] Error updating websocket:", err)
		return
	}
	defer conn.Close()

	p.mu.Lock()
	p.clients[conn] = true
	p.mu.Unlock()

	log.Println("[Listener] New client connected")

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			p.mu.Lock()
			delete(p.clients, conn)
			p.mu.Unlock()
			fmt.Println("[Listener] Client disconnected")
			break
		}
	}
}
