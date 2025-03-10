package usecases

import "net/http"

type WSBroadcasterAdapter interface {
	ServeWebsocketBroadcasterHTTP(w http.ResponseWriter, r *http.Request)
}
