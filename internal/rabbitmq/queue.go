package rabbitmq

import (
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	rabbitmq "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	RetryTime    = 3 * time.Second
	QueueNameTag = "queueName"
	MessageTag   = "message"
	ErrorTag     = "error"
)

type Queue interface {
	SendMessage([]byte) error
	ReadMessage() ([]byte, error)
	Close() error
}

type RabbitMQQueue struct {
	queue   rabbitmq.Queue
	channel *rabbitmq.Channel
	logger  *log.Logger
}

var _ Queue = (*RabbitMQQueue)(nil)

type QueueConfig struct {
	Port   *int
	Name   string
	Logger *log.Logger
}

func NewQueueConnection(config QueueConfig) (*RabbitMQQueue, error) {
	if config.Logger == nil {
		return nil, errors.New("option 'Logger' is mandatory")
	}
	if config.Port == nil {
		return nil, errors.New("option 'Port' is mandatory")
	}
	if len(config.Name) == 0 {
		config.Name = uuid.Must(uuid.NewV4()).String()
	}

	url := fmt.Sprintf("amqp://guest:guest@localhost:%d/", *config.Port)
	conn, err := rabbitmq.Dial(url)
	if err != nil {
		config.Logger.WithFields(log.Fields{
			ErrorTag: err.Error(),
		}).Errorf("[RabbitMQQueue] error connecting with RabbitMQ")
		return nil, ErrConnecting
	}

	ch, err := conn.Channel()
	if err != nil {
		config.Logger.WithFields(log.Fields{
			ErrorTag: err.Error(),
		}).Error("[RabbitMQQueue] error opening RabbitMQ channel")
		return nil, ErrOpening
	}

	q, err := ch.QueueDeclare(config.Name, false, false, false, false, nil)
	if err != nil {
		config.Logger.WithFields(log.Fields{
			ErrorTag: err.Error(),
		}).Error("[RabbitMQQueue] error in queue declaration")
		return nil, ErrQueueDeclaration
	}

	config.Logger.WithFields(log.Fields{
		QueueNameTag: q.Name,
	}).Info("[RabbitMQQueue] queue successfully created")

	return &RabbitMQQueue{
		logger:  config.Logger,
		queue:   q,
		channel: ch,
	}, nil
}

func (q RabbitMQQueue) SendMessage(message []byte) error {
	q.logger.WithFields(log.Fields{
		QueueNameTag: q.queue.Name,
		MessageTag:   string(message),
	}).Info("[RabbitMQQueue] sending message")

	queueMessage := rabbitmq.Publishing{
		ContentType: "text/plain",
		Body:        message,
	}

	err := q.channel.Publish("", q.queue.Name, false, false, queueMessage)
	if err != nil {
		q.logger.WithFields(log.Fields{
			QueueNameTag: q.queue.Name,
			MessageTag:   string(message),
			ErrorTag:     err.Error(),
		}).Error("[RabbitMQQueue] error sending message to RabbitMQ")
		return ErrSendingMessage
	}

	q.logger.WithFields(log.Fields{
		QueueNameTag: q.queue.Name,
		MessageTag:   string(message),
	}).Debug("[RabbitMQQueue] message sent")
	return nil
}

func (q RabbitMQQueue) ReadMessage() ([]byte, error) {
	q.logger.WithFields(log.Fields{
		QueueNameTag: q.queue.Name,
	}).Info("[RabbitMQQueue] looking for messages")
	for {
		delivery, ok, err := q.channel.Get(q.queue.Name, true)
		if err != nil {
			if errors.Is(err, rabbitmq.ErrClosed) {
				q.logger.WithFields(log.Fields{
					QueueNameTag: q.queue.Name,
					ErrorTag:     err.Error(),
				}).Error("[RabbitMQQueue] error reading message because channel is closed")
				return nil, ErrChannelClosed
			}

			q.logger.WithFields(log.Fields{
				QueueNameTag: q.queue.Name,
				ErrorTag:     err.Error(),
			}).Warnf("[RabbitMQQueue] error reading message, trying again in %d seconds", RetryTime)
			time.Sleep(RetryTime)
		}

		if !ok {
			continue
		}

		q.logger.WithFields(log.Fields{
			QueueNameTag: q.queue.Name,
		}).Debug("[RabbitMQQueue] message read")

		return delivery.Body, nil
	}

}

func (q RabbitMQQueue) Close() error {
	q.logger.WithFields(log.Fields{
		QueueNameTag: q.queue.Name,
	}).Info("[RabbitMQQueue] closing channel message")

	err := q.channel.Close()
	if err != nil {
		q.logger.WithFields(log.Fields{
			QueueNameTag: q.queue.Name,
			ErrorTag:     err.Error(),
		}).Error("[RabbitMQQueue] error closing channel")
		return ErrClosingChannel
	}

	q.logger.WithFields(log.Fields{
		QueueNameTag: q.queue.Name,
	}).Debug("[RabbitMQQueue] channel closed")

	return nil
}

var ErrConnecting = errors.New("error connecting with RabbitMQ")
var ErrOpening = errors.New("error opening RabbitMQ channel")
var ErrQueueDeclaration = errors.New("error in queue declaration")
var ErrSendingMessage = errors.New("error sending message to RabbitMQ")
var ErrChannelClosed = errors.New("error reading message because channel is closed")
var ErrClosingChannel = errors.New("error closing channel")
