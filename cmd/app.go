package main

import (
	"gig-assessment/pkg/microservices"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type App struct {
	WSReceiver    *microservices.WSReceiver
	WSBroadcaster *microservices.WSBroadcaster
}

func (a App) serveWebsocketReceiverHTTP(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	} //todo: see the default values

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}
	//defer conn.Close()

	go a.WSReceiver.ReadMessage(conn)
}

func (a App) serveWebsocketBroadcasterHTTP(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("[WSBroadcaster] Error updating websocket:", err)
		return
	}
	//defer conn.Close()

	go a.WSBroadcaster.ReadClientMessage(conn)
}

func (a App) startBroadcaster() {
	_, err := a.WSBroadcaster.ReadQueueMessage()
	if err != nil {
		panic(err)
	}
}

func (a App) startWebsocketReceiverServer() {
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func (a App) startWebsocketBroadcasterServer() {
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		panic(err)
	}
}
