package main

import (
	"gig-assessment/internal/rabbitmq-queue"
	"gig-assessment/pkg/wsbroadcaster"
	"gig-assessment/pkg/wsreceiver"
	"net/http"
)

func main() {
	var err error

	var config rabbitmqqueue.QueueConfig
	config.Host = "localhost"
	config.Port = 5672
	queue, err := rabbitmqqueue.NewQueueConnection(config)
	if err != nil {
		panic(err)
	}

	var newWSReceiverOptions wsreceiver.NewWSReceiverOptions
	newWSReceiverOptions.Queue = queue
	receiver, err := wsreceiver.NewWSReceiver(newWSReceiverOptions)
	if err != nil {
		panic(err)
	}

	var newWSBroadcasterOptions wsbroadcaster.NewWSBroadcasterOptions
	newWSBroadcasterOptions.Queue = queue
	broadcaster, err := wsbroadcaster.NewWSBroadcaster(newWSBroadcasterOptions)
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
