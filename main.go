package main

import (
	"gig-assessment/listener"
	"gig-assessment/publisher"
	"gig-assessment/rabbitmq-queue"
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

	var newPublisherOptions publisher.NewPublisherOptions
	newPublisherOptions.Queue = queue
	p, err := publisher.NewPublisher(newPublisherOptions)
	if err != nil {
		panic(err)
	}

	var newListenerOptions listener.NewListenerOptions
	newListenerOptions.Queue = queue
	l, err := listener.NewListener(newListenerOptions)
	if err != nil {
		panic(err)
	}

	go func() {
		l.ReadMessage()
	}()

	// register handlers
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		p.ManageMessage(w, r)
	})

	http.ListenAndServe(":8080", nil)

}
