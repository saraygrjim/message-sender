package wsreceiver

import (
	"bytes"
	"errors"
	"fmt"
	"gig-assessment/internal/rabbitmq-queue"
	"github.com/gorilla/websocket"
	"log"
)

type WSReceiver struct {
	queue *rabbitmqqueue.RabbitMQQueue
}

type NewWSReceiverOptions struct {
	Queue *rabbitmqqueue.RabbitMQQueue
}

func NewWSReceiver(opts NewWSReceiverOptions) (*WSReceiver, error) {
	if opts.Queue == nil {
		return nil, errors.New("option 'RabbitMQQueue' is mandatory")
	}

	return &WSReceiver{
		queue: opts.Queue,
	}, nil
}

func (wsr *WSReceiver) ReadMessage(conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("[WSReceiver] there was an error reading the message: %s \n", err.Error())
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				panic(err) //todo: this doesnt work, review it
			}
		}

		go wsr.SendMessageToQueue(message)
	}
}

func (wsr *WSReceiver) SendMessageToQueue(message []byte) {
	//todo: I am assuming that this is parsing the message, improve it with a function
	message = bytes.TrimSpace(bytes.Replace(message, []byte("\n"), []byte(" "), -1))
	fmt.Println(string(message))

	err := wsr.queue.SendMessage(message)
	if err != nil {
		log.Fatalf("[WSReceiver] there was an error writing message '%s' into queue", message)
	}
}
