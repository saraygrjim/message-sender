package usecases_test

import (
	"errors"
	"github.com/gofrs/uuid"
	infrawebsocket "github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"message-sender/internal/rabbitmq"
	"message-sender/internal/websocket"
	"message-sender/microservices/receiver/pkg/usecases"
	"testing"
)

var _ rabbitmq.Queue = (*MockQueue)(nil)

type MockQueue struct {
	mock.Mock
}

func (m *MockQueue) SendMessage(msg []byte) error {
	args := m.Called(msg)
	return args.Error(0)
}

func (m *MockQueue) ReadMessage() ([]byte, error) {
	args := m.Called()
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockQueue) Close() error {
	args := m.Called()
	return args.Error(0)
}

var _ websocket.Connection = (*MockWebSocketConnection)(nil)

type MockWebSocketConnection struct {
	mock.Mock
}

func (m *MockWebSocketConnection) WriteMessage(_ []byte) error {
	return errors.New("not implemented")
}

func (m *MockWebSocketConnection) RemoteAddr() string {
	return ""
}

func (m *MockWebSocketConnection) ReadMessage() ([]byte, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockWebSocketConnection) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestNewWSReceiver_Ok(t *testing.T) {
	logger := logrus.New()
	queue := new(MockQueue)

	t.Run("Should be able to create a receiver", func(t *testing.T) {
		receiver, err := usecases.NewWSReceiver(usecases.DefaultWSReceiverOptions{
			Queue:  queue,
			Logger: logger,
		})
		assert.NoError(t, err)
		assert.NotNil(t, receiver)
	})
}

func TestNewWSReceiver_Error(t *testing.T) {
	logger := logrus.New()
	queue := new(MockQueue)

	t.Run("Should not be able to create a receiver if Queue is equal to nil", func(t *testing.T) {
		_, err := usecases.NewWSReceiver(usecases.DefaultWSReceiverOptions{
			Queue:  nil,
			Logger: logger,
		})
		require.Error(t, err)
		require.Equal(t, "option 'Queue' is mandatory", err.Error())
	})

	t.Run("Should no be able to crate a receiver if Logger is equal to nil", func(t *testing.T) {
		_, err := usecases.NewWSReceiver(usecases.DefaultWSReceiverOptions{
			Queue:  queue,
			Logger: nil,
		})
		require.Error(t, err)
		require.Equal(t, "option 'Logger' is mandatory", err.Error())
	})
}

func TestSendMessageToQueue_Ok(t *testing.T) {
	logger := logrus.New()
	queue := new(MockQueue)

	receiver, err := usecases.NewWSReceiver(usecases.DefaultWSReceiverOptions{
		Queue:  queue,
		Logger: logger,
	})
	require.NoError(t, err)

	// precondition
	queue.On("SendMessage", mock.Anything).Return(nil)

	message := []byte(uuid.Must(uuid.NewV4()).String())
	err = receiver.SendMessageToQueue(message)
	require.Nil(t, err)
	queue.AssertCalled(t, "SendMessage", message)
}

func TestSendMessageToQueue_Error(t *testing.T) {
	logger := logrus.New()
	queue := new(MockQueue)

	receiver, err := usecases.NewWSReceiver(usecases.DefaultWSReceiverOptions{
		Queue:  queue,
		Logger: logger,
	})
	require.NoError(t, err)

	// precondition
	queue.On("SendMessage", mock.Anything).Return(errors.New("Queue failure"))

	message := []byte(uuid.Must(uuid.NewV4()).String())
	err = receiver.SendMessageToQueue(message)
	require.NotNil(t, err)
	require.True(t, errors.Is(err, usecases.ErrWritingMessageInQueue))

	queue.AssertCalled(t, "SendMessage", message)
}

func TestReadMessage_Ok(t *testing.T) {
	logger := logrus.New()
	queue := new(MockQueue)
	conn := new(MockWebSocketConnection)

	receiver, err := usecases.NewWSReceiver(usecases.DefaultWSReceiverOptions{
		Queue:  queue,
		Logger: logger,
	})
	assert.NoError(t, err)

	message := []byte(uuid.Must(uuid.NewV4()).String())

	// preconditions
	conn.On("ReadMessage").Return(message, nil).Once()
	conn.On("ReadMessage").Return(nil, nil).Once()
	conn.On("ReadMessage").Return(nil, &infrawebsocket.CloseError{
		Code: infrawebsocket.CloseGoingAway,
		Text: "bad close",
	}).Once()
	conn.On("Close").Return(nil)
	queue.On("SendMessage", mock.Anything).Return(nil)

	// test
	err = receiver.ReadMessage(conn)
	require.Nil(t, err)

	// message is sent to queue
	queue.AssertCalled(t, "SendMessage", message)

	// connection is closed
	conn.AssertCalled(t, "Close")

}

func TestReadMessage_Error(t *testing.T) {
	logger := logrus.New()
	queue := new(MockQueue)
	conn := new(MockWebSocketConnection)

	receiver, err := usecases.NewWSReceiver(usecases.DefaultWSReceiverOptions{
		Queue:  queue,
		Logger: logger,
	})
	assert.NoError(t, err)

	t.Run("Should manage the error when reading messages", func(t *testing.T) {
		// precondition
		conn.On("ReadMessage").Return(nil, errors.New("read error")).Once()
		conn.On("Close").Return(nil)

		err = receiver.ReadMessage(conn)
		require.NotNil(t, err)
		require.True(t, errors.Is(err, usecases.ErrReadingWSMessage))

		// no messages sent
		queue.AssertNotCalled(t, "SendMessage", []byte{})

		// connection is closed
		conn.AssertCalled(t, "Close")
	})

	t.Run("Should manage the error wen websocket unexpectedly closes", func(t *testing.T) {
		// precondition
		conn.On("ReadMessage").Return(nil,
			&infrawebsocket.CloseError{
				Code: infrawebsocket.CloseGoingAway,
				Text: "bad close",
			}).Once()
		conn.On("Close").Return(nil)

		err = receiver.ReadMessage(conn)
		require.Nil(t, err)

		// no messages sent
		queue.AssertNotCalled(t, "SendMessage", []byte{})

		// connection is closed
		conn.AssertCalled(t, "Close")
	})
}
