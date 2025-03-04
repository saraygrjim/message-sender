package main

import (
	"gig-assessment/publisher"
	"net/http"
)

func main() {
	var err error

	var newPublisherOptions publisher.NewPublisherOptions
	newPublisherOptions.QueueHost = "localhost"
	newPublisherOptions.QueuePort = 5672
	p, err := publisher.NewPublisher(newPublisherOptions)
	if err != nil {
		panic(err)
	}

	// register handlers
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		p.ManageMessage(w, r)
	})

	http.ListenAndServe(":8080", nil)

}
