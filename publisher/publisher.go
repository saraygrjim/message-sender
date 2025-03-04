package publisher

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
)

type Publisher struct {
	queue amqp.Queue
	channel *amqp.Channel
	ws    *websocket.Conn
}

type NewPublisherOptions struct {
	QueueHost string
	QueuePort int
}


func NewPublisher(opts NewPublisherOptions) (*Publisher, error) {
	if len(opts.QueueHost) == 0 {
		opts.QueueHost = "localhost"
	}

	if opts.QueuePort == 0 {
		opts.QueuePort = 5672
	}

	url := fmt.Sprintf("amqp://guest:guest@%s:%d/", opts.QueueHost, opts.QueuePort)
	fmt.Println(url)
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatal("[Publisher] error connecting with RabbitMQ: ", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("[Publisher] error opening RabbitMQ channel: ", err)
	}

	q, err := ch.QueueDeclare("messages", false, false, false, false, nil)
	if err != nil {
		log.Fatal("[Publisher] error in queue declaration:", err)
	}

	return &Publisher{
		queue: q,
		channel: ch,
		ws:    nil,
	}, nil
}

func (p *Publisher) InitWebSocketConnection(w http.ResponseWriter, r *http.Request) error {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	} //todo: see the default values

	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	//todo: review this
	connection.SetReadLimit(512)
	connection.SetReadDeadline(time.Now().Add(60 * time.Second))
	connection.SetPongHandler(func(string) error { connection.SetReadDeadline(time.Now().Add(60 * time.Second)); return nil })

	p.ws = connection

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
			if  websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				panic(errors.New(msg)) //todo: this doesnt work, review it
			}
			return errors.New(msg)
		}

		//todo: I am assuming that this is parsing the message, improve it with a function
		message = bytes.TrimSpace(bytes.Replace(message, []byte("\n"), []byte(" "), -1))
		fmt.Println(string(message))

		err = p.SendMessage(message)
		if err != nil {
			return err
		}
	}

}

//todo: create a queue package to have all of this functions centralized
func (p *Publisher) SendMessage(message []byte) error {
	queueMessage :=  amqp.Publishing{
		ContentType: "text/plain",
		Body: message,
	}

	err := p.channel.Publish("", p.queue.Name, false, false, queueMessage)
	if err != nil {
		msg := fmt.Sprintf("[Publisher] error sending message to RabbitMQ: %s", err.Error())
		log.Print(msg)
		return errors.New(msg)
	}

	return nil
}
