package listener

import (
	"gig-assessment/rabbitmq-queue"
	"github.com/gorilla/websocket"
	"errors"
)

type Listener struct {
	queue *rabbitmqqueue.Queue
	ws    *websocket.Conn
}

type NewListenerOptions struct {
	Queue *rabbitmqqueue.Queue

}

func NewListener(opts NewListenerOptions) (*Listener, error) {
	if opts.Queue == nil {
		return nil, errors.New("requiered 'Queue' equal to nil")
	}

	return &Listener{
		queue: opts.Queue,
		ws:    nil,
	}, nil
}




func (p Listener) ReadMessage() error {
	for {
		p.queue.ReadMessage()
	}
}