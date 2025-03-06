package rabbitmqqueue

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"log"
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
		log.Fatal("[Queue] error connecting with RabbitMQ: ", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("[Queue] error opening RabbitMQ channel: ", err)
	}

	q, err := ch.QueueDeclare("messages", false, false, false, false, nil)
	if err != nil {
		log.Fatal("[Queue] error in queue declaration:", err)
	}

	return &Queue{
		queue:   q,
		channel: ch,
	}, nil
}

func (q Queue) SendMessage(message []byte) error {
	queueMessage := amqp.Publishing{
		ContentType: "text/plain",
		Body:        message,
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
	for {
		delivery, ok, err := q.channel.Get(q.queue.Name, true)
		if err != nil {
			log.Fatal("[Queue] error reading messages:", err)
		}

		if !ok {
			continue
		}

		fmt.Printf("[Queue] %s\n", string(delivery.Body))
		return delivery.Body, nil
	}

}
