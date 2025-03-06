package rabbitmqqueue

import (
	"errors"
	"fmt"
	rabbitmq "github.com/streadway/amqp"
	"log"
)

type RabbitMQQueue struct {
	queue   rabbitmq.Queue
	channel *rabbitmq.Channel
}

type Queue interface {
	SendMessage([]byte) error
	ReadMessage() ([]byte, error)
}

var _ Queue = (*RabbitMQQueue)(nil)

type QueueConfig struct {
	Host string
	Port int
}

func NewQueueConnection(config QueueConfig) (*RabbitMQQueue, error) {
	if len(config.Host) == 0 {
		config.Host = "localhost"
	}

	if config.Port == 0 {
		config.Port = 5672
	}

	url := fmt.Sprintf("amqp://guest:guest@%s:%d/", config.Host, config.Port)
	fmt.Println(url)
	conn, err := rabbitmq.Dial(url)
	if err != nil {
		log.Fatal("[RabbitMQQueue] error connecting with RabbitMQ: ", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("[RabbitMQQueue] error opening RabbitMQ channel: ", err)
	}

	q, err := ch.QueueDeclare("messages", false, false, false, false, nil)
	if err != nil {
		log.Fatal("[RabbitMQQueue] error in queue declaration:", err)
	}

	return &RabbitMQQueue{
		queue:   q,
		channel: ch,
	}, nil
}

func (q RabbitMQQueue) SendMessage(message []byte) error {
	queueMessage := rabbitmq.Publishing{
		ContentType: "text/plain",
		Body:        message,
	}

	err := q.channel.Publish("", q.queue.Name, false, false, queueMessage)
	if err != nil {
		msg := fmt.Sprintf("[RabbitMQQueue] error sending message to RabbitMQ: %s", err.Error())
		log.Print(msg)
		return errors.New(msg)
	}

	return nil
}

func (q RabbitMQQueue) ReadMessage() ([]byte, error) {
	for {
		delivery, ok, err := q.channel.Get(q.queue.Name, true)
		if err != nil {
			log.Printf("[RabbitMQQueue] error reading messages: %s", err)
			if errors.Is(err, rabbitmq.ErrClosed) {
				return nil, err
			}
		}

		if !ok {
			continue
		}

		return delivery.Body, nil
	}

}

func (q RabbitMQQueue) Close() error {
	return q.channel.Close()
}
