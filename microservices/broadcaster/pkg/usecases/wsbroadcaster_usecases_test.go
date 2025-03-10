package usecases_test

import (
	"github.com/go-errors/errors"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"message-sender/internal/rabbitmq"
	"message-sender/internal/websocket"
	"message-sender/microservices/broadcaster/pkg/usecases"
	"testing"
)

var _ rabbitmq.Queue = (*MockQueue)(nil)

type MockQueue struct {
	mock.Mock
}

func (m *MockQueue) SendMessage(_ []byte) error {
	return errors.New("not implemented")
}

func (m *MockQueue) ReadMessage() ([]byte, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockQueue) Close() error {
	args := m.Called()
	return args.Error(0)
}

var _ websocket.Connection = (*MockWSConnection)(nil)

type MockWSConnection struct {
	mock.Mock
}

func (m *MockWSConnection) ReadMessage() ([]byte, error) {
	args := m.Called()
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockWSConnection) WriteMessage(message []byte) error {
	args := m.Called(message)
	return args.Error(0)
}

func (m *MockWSConnection) RemoteAddr() string {
	return "mock-address"
}

func (m *MockWSConnection) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestNewWSBroadcaster_Success(t *testing.T) {
	logger := log.New()
	queue := new(MockQueue)

	var opts = usecases.DefaultWSBroadcasterOptions{
		Queue:  queue,
		Logger: logger,
	}

	broadcaster, err := usecases.NewWSBroadcaster(opts)
	assert.NoError(t, err)
	assert.NotNil(t, broadcaster)
}

func TestNewWSBroadcaster_Error(t *testing.T) {
	logger := log.New()
	queue := new(MockQueue)

	t.Run("Should not be able to create a broadcaster if Queue is equal to nil", func(t *testing.T) {
		var opts = usecases.DefaultWSBroadcasterOptions{
			Queue:  nil,
			Logger: logger,
		}
		_, err := usecases.NewWSBroadcaster(opts)
		require.Error(t, err)
		require.Equal(t, "option 'Queue' is mandatory", err.Error())
	})

	t.Run("Should no be able to crate a broadcaster if Logger is equal to nil", func(t *testing.T) {
		var opts = usecases.DefaultWSBroadcasterOptions{
			Queue:  queue,
			Logger: nil,
		}
		_, err := usecases.NewWSBroadcaster(opts)
		require.Error(t, err)
		require.Equal(t, "option 'Logger' is mandatory", err.Error())
	})
}

func TestReadAndSend_Success(t *testing.T) {
	logger := log.New()
	queue := new(MockQueue)
	var opts = usecases.DefaultWSBroadcasterOptions{
		Queue:  queue,
		Logger: logger,
	}
	broadcaster, err := usecases.NewWSBroadcaster(opts)
	require.NoError(t, err)

	msg := []byte(uuid.Must(uuid.NewV4()).String())

	// precondition
	queue.On("ReadMessage").Return(msg, nil).Once()
	queue.On("ReadMessage").Return(nil, rabbitmq.ErrChannelClosed).Once()

	err = broadcaster.ReadAndSend()
	require.NotNil(t, err)
}

func TestSendMessageToWS_Success(t *testing.T) {
	queue := new(MockQueue)
	conn := new(MockWSConnection)
	logger := log.New()
	var opts = usecases.DefaultWSBroadcasterOptions{
		Queue:  queue,
		Logger: logger,
	}
	broadcaster, err := usecases.NewWSBroadcaster(opts)
	assert.NotNil(t, broadcaster)
	assert.Nil(t, err)

	broadcaster.NewClient(conn)

	msg := []byte(uuid.Must(uuid.NewV4()).String())
	conn.On("WriteMessage", msg).Return(nil)

	err = broadcaster.SendMessageToWS(msg)
	assert.NoError(t, err)
	conn.AssertExpectations(t)
}

func TestSendMessageToWS_Failure(t *testing.T) {
	queue := new(MockQueue)
	conn := new(MockWSConnection)
	logger := log.New()
	var opts = usecases.DefaultWSBroadcasterOptions{
		Queue:  queue,
		Logger: logger,
	}
	broadcaster, err := usecases.NewWSBroadcaster(opts)
	assert.NotNil(t, broadcaster)
	assert.Nil(t, err)

	broadcaster.NewClient(conn)

	msg := []byte(uuid.Must(uuid.NewV4()).String())
	conn.On("WriteMessage", msg).Return(errors.New("write error"))

	err = broadcaster.SendMessageToWS(msg)
	assert.Error(t, err)
	assert.Equal(t, usecases.ErrSendingMessageToWebsocket, err)
	conn.AssertExpectations(t)
}
