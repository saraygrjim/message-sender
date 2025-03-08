package websocket

import (
	"errors"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const (
	ErrorTag = "error"
)

type Connection struct {
	conn   *websocket.Conn
	logger *log.Logger
}

func NewConnection(w http.ResponseWriter, r *http.Request, logger *log.Logger) (*Connection, error) {
	if logger == nil {
		return nil, errors.New("option 'Logger' is mandatory")
	}

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.WithFields(log.Fields{
			ErrorTag: err.Error(),
		}).Fatal("[Websocket] Error updating websocket")
		return nil, err
	}

	return &Connection{
		logger: logger,
		conn:   conn,
	}, nil
}

func NewClient(url string, logger *log.Logger) (*Connection, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		logger.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("Error connecting with websocket")
		return nil, err
	}

	return &Connection{
		logger: logger,
		conn:   conn,
	}, nil
}

func (c Connection) Close() error {
	return c.conn.Close()
}

func (c Connection) ReadMessage() ([]byte, error) {
	_, msg, err := c.conn.ReadMessage()
	return msg, err
}

func (c Connection) WriteMessage(msg []byte) error {
	return c.conn.WriteMessage(websocket.TextMessage, msg)
}

func (c Connection) RemoteAddr() string {
	return c.conn.RemoteAddr().String()
}

func IsUnexpectedCloseError(err error) bool {
	return websocket.IsUnexpectedCloseError(err)
}
