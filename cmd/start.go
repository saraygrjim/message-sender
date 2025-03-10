package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"message-sender/internal/rabbitmq"
	broadcasterGraph "message-sender/microservices/broadcaster/pkg/graph"
	receiverGraph "message-sender/microservices/receiver/pkg/graph"
	subscriberGraph "message-sender/microservices/subscriber/pkg/graph"
	"net/http"
	"os"
	"time"
)

func Receiver() {
	var err error

	logger := log.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&log.TextFormatter{
		TimestampFormat: time.DateTime,
		FullTimestamp:   true,
	})

	port := 5672
	var config rabbitmq.QueueConfig
	config.Port = &port
	config.Logger = logger
	config.Name = "messages-sender-queue"
	queue, err := rabbitmq.NewQueueConnection(config)
	if err != nil {
		panic(err)
	}

	var receiverOpts receiverGraph.ReceiverOptions
	receiverOpts.Logger = logger
	receiverOpts.Queue = queue
	receiver, err := receiverGraph.Install(receiverOpts)
	if err != nil {
		panic(err)
	}

	go receiver.StartWebsocketReceiverServer()
	serveHTML(logger)
}

func serveHTML(l *log.Logger) {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	port := 9090
	l.Infof("Server running on http://localhost:%d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		l.WithFields(log.Fields{"error": err}).Error("Error starting server")
		panic(err)
	}
}

func Broadcaster() {
	var err error

	logger := log.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&log.TextFormatter{
		TimestampFormat: time.DateTime,
		FullTimestamp:   true,
	})

	port := 5672
	var config rabbitmq.QueueConfig
	config.Port = &port
	config.Logger = logger
	config.Name = "messages-sender-queue"
	queue, err := rabbitmq.NewQueueConnection(config)
	if err != nil {
		panic(err)
	}

	var broadcasterOpts broadcasterGraph.BroadcasterOptions
	broadcasterOpts.Logger = logger
	broadcasterOpts.Queue = queue
	broadcaster, err := broadcasterGraph.Install(broadcasterOpts)
	if err != nil {
		panic(err)
	}

	go broadcaster.StartWebsocketBroadcasterServer()

	broadcaster.StartBroadcaster()
}

func StartSubscriber() {
	logger := log.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&log.TextFormatter{
		TimestampFormat: time.DateTime,
		FullTimestamp:   true,
	})

	var subscriberOpts subscriberGraph.SubscriberOptions
	subscriberOpts.Logger = logger
	subscriber, err := subscriberGraph.Install(subscriberOpts)
	if err != nil {
		panic(err)
	}

	err = subscriber.StartSubscriber()
	if err != nil {
		panic(err)
	}
}
