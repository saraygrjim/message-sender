package rabbitmqqueue

import (
	"errors"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type Queue struct {
	queue   amqp.Queue
	channel *amqp.Channel
}

type QueueConfig struct {
	Host string
	Port int
}

func NewQueueConnection(config QueueConfig) (*Queue, error) {
	if len(config.Host) == 0 {
		config.Host = "localhost"
	}

	if config.Port == 0 {
		config.Port = 5672
	}

	url := fmt.Sprintf("amqp://guest:guest@%s:%d/", config.Host, config.Port)
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

	return &Queue{
		queue: q,
		channel: ch,
	}, nil
}


func (q Queue) SendMessage(message []byte) error {
	queueMessage :=  amqp.Publishing{
		ContentType: "text/plain",
		Body: message,
	}

	err := q.channel.Publish("", q.queue.Name, false, false, queueMessage)
	if err != nil {
		msg := fmt.Sprintf("[Queue] error sending message to RabbitMQ: %s", err.Error())
		log.Print(msg)
		return errors.New(msg)
	}

	return nil
}

func (q Queue) ReadMessage() ([]byte, error) {
	msgs, err := q.channel.Consume(q.queue.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal("Error al consumir mensajes:", err)
	}

	go func() {
		for msg := range msgs {
			fmt.Printf("[Listener] %s\n", string(msg.Body))
		}
	}()

	return nil, nil
}