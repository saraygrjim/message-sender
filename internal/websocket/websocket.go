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

type Connection interface {
	Close() error
	ReadMessage() ([]byte, error)
	WriteMessage(msg []byte) error
	RemoteAddr() string
}

func NewConnection(w http.ResponseWriter, r *http.Request, logger *log.Logger) (*WSConnection, error) {
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

	return &WSConnection{
		logger: logger,
		conn:   conn,
	}, nil
}

func NewClient(url string, logger *log.Logger) (*WSConnection, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		logger.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("Error connecting with websocket")
		return nil, err
	}

	return &WSConnection{
		logger: logger,
		conn:   conn,
	}, nil
}

type WSConnection struct {
	conn   *websocket.Conn
	logger *log.Logger
}

var _ Connection = (*WSConnection)(nil)

func (c WSConnection) Close() error {
	return c.conn.Close()
}

func (c WSConnection) ReadMessage() ([]byte, error) {
	_, msg, err := c.conn.ReadMessage()
	return msg, err
}

func (c WSConnection) WriteMessage(msg []byte) error {
	return c.conn.WriteMessage(websocket.TextMessage, msg)
}

func (c WSConnection) RemoteAddr() string {
	return c.conn.RemoteAddr().String()
}

func IsUnexpectedCloseError(err error) bool {
	return websocket.IsUnexpectedCloseError(err)
}
