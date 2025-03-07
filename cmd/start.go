package main

import (
	"gig-assessment/internal/rabbitmq-queue"
	"gig-assessment/pkg/microservices"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

func main() {
	var err error

	logger := log.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(log.InfoLevel) // TODO: set as argument
	logger.SetFormatter(&log.TextFormatter{
		TimestampFormat: time.DateTime,
		FullTimestamp:   true,
	})

	var config rabbitmqqueue.QueueConfig // TODO: set config as arguments
	config.Host = "localhost"
	config.Port = 5672
	config.Logger = logger
	queue, err := rabbitmqqueue.NewQueueConnection(config)
	if err != nil {
		panic(err)
	}

	var newWSReceiverOptions microservices.NewWSReceiverOptions
	newWSReceiverOptions.Queue = queue
	newWSReceiverOptions.Logger = logger
	receiver, err := microservices.NewWSReceiver(newWSReceiverOptions)
	if err != nil {
		panic(err)
	}

	var newWSBroadcasterOptions microservices.NewWSBroadcasterOptions
	newWSBroadcasterOptions.Queue = queue
	newWSBroadcasterOptions.Logger = logger
	broadcaster, err := microservices.NewWSBroadcaster(newWSBroadcasterOptions)
	if err != nil {
		panic(err)
	}

	app := App{
		WSReceiver:    receiver,
		WSBroadcaster: broadcaster,
	}

	// register handlers
	http.HandleFunc("/ws", app.serveWebsocketReceiverHTTP)

	http.HandleFunc("/echo", app.serveWebsocketBroadcasterHTTP)

	go app.startWebsocketReceiverServer()

	go app.startWebsocketBroadcasterServer()

	app.startBroadcaster()
}
