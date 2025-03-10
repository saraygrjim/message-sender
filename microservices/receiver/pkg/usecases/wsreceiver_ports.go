package usecases

import "net/http"

type WSReceiverAdapter interface {
	ServeWebsocketReceiverHTTP(w http.ResponseWriter, r *http.Request)
}
