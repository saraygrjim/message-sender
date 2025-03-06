package publisher

import (
	"bytes"
	"errors"
	"fmt"
	"gig-assessment/rabbitmq-queue"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

type Publisher struct {
	queue *rabbitmqqueue.Queue
	ws    *websocket.Conn
}

type NewPublisherOptions struct {
	Queue *rabbitmqqueue.Queue
}

func NewPublisher(opts NewPublisherOptions) (*Publisher, error) {
	if opts.Queue == nil {
		return nil, errors.New("option 'Queue' is mandatory")
	}

	return &Publisher{
		queue: opts.Queue,
		ws:    nil,
	}, nil
}

func (p *Publisher) InitWebSocketConnection(w http.ResponseWriter, r *http.Request) error {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	} //todo: see the default values

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	//todo: review this parameters
	conn.SetReadLimit(512)
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(60 * time.Second)); return nil })

	p.ws = conn

	return nil
}

// ManageMessage reads messages from the websocket and sends them to the queue
func (p *Publisher) ManageMessage(w http.ResponseWriter, r *http.Request) error {
	err := p.InitWebSocketConnection(w, r)
	if err != nil {
		msg := fmt.Sprintf("there was an error trying to init the websocket connection: %s ", err.Error())
		log.Printf("[Publisher] %s\n", msg)
		return errors.New(msg)
	}

	for {
		_, message, err := p.ws.ReadMessage()
		if err != nil {
			msg := fmt.Sprintf("there was an error reading the message: %s ", err.Error())
			log.Printf("[Publisher] %s\n", msg)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				panic(errors.New(msg)) //todo: this doesnt work, review it
			}
			return errors.New(msg)
		}

		//todo: I am assuming that this is parsing the message, improve it with a function
		message = bytes.TrimSpace(bytes.Replace(message, []byte("\n"), []byte(" "), -1))
		fmt.Println(string(message))

		err = p.queue.SendMessage(message)
		if err != nil {
			return err
		}
	}

}
