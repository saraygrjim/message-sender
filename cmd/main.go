package main

import (
	"gig-assessment/internal/rabbitmq-queue"
	"gig-assessment/pkg/usecases/listener"
	"gig-assessment/pkg/usecases/publisher"
	"net/http"
)

func main() {
	var err error

	var config rabbitmqqueue.rabbitmqqueue
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

	// register handlers
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		err = p.ManageMessage(w, r)
		if err != nil {
			panic(err)
		}
	})

	http.HandleFunc("/echo", l.HandleWebSocket) // Maneja conexiones WebSocket

	go func() {
		err = http.ListenAndServe(":8080", nil)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		err = http.ListenAndServe(":8081", nil)
		if err != nil {
			panic(err)
		}
	}()

	//err = l.InitWebSocketConnection()
	//if err != nil {
	//	panic(err)
	//}

	_, err = l.ReadMessage()
	if err != nil {
		panic(err)
	}

}
