package rabbitmqqueue

import (
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewQueueConnection(t *testing.T) {
	t.Run("Should be able to connect a queue without Host", func(t *testing.T) {
		config := QueueConfig{Port: 5672}

		queue, err := NewQueueConnection(config)
		require.NotNil(t, queue)
		require.Nil(t, err)
		require.Equal(t, "messages", queue.queue.Name)
	})

	t.Run("Should be able to connect a queue without Port", func(t *testing.T) {
		config := QueueConfig{Host: "localhost"}

		queue, err := NewQueueConnection(config)
		require.NotNil(t, queue)
		require.Nil(t, err)
		require.Equal(t, "messages", queue.queue.Name)
	})

	t.Run("Should be able to connect a queue with all the parameters", func(t *testing.T) {

		config := QueueConfig{Host: "localhost", Port: 5672}

		queue, err := NewQueueConnection(config)
		require.NotNil(t, queue)
		require.Nil(t, err)
		require.Equal(t, "messages", queue.queue.Name)
	})
}

func TestSendMessage_Success(t *testing.T) {
	config := QueueConfig{Host: "localhost", Port: 5672}
	queue, err := NewQueueConnection(config)
	assert.Nil(t, err)

	message := uuid.Must(uuid.NewV4()).String()
	err = queue.SendMessage([]byte(message))
	require.Nil(t, err)

	defer func() {
		_, err = queue.ReadMessage()
		assert.Nil(t, err)

		err = queue.Close()
		assert.Nil(t, err)
	}()
}

func TestSendMessage_Error(t *testing.T) {
	config := QueueConfig{Host: "localhost", Port: 5672}
	queue, err := NewQueueConnection(config)
	assert.Nil(t, err)

	err = queue.Close()
	assert.Nil(t, err)

	message := uuid.Must(uuid.NewV4()).String()
	err = queue.SendMessage([]byte(message))
	require.NotNil(t, err)
}

func TestReadMessage_Success(t *testing.T) {
	config := QueueConfig{Host: "localhost", Port: 5672}
	queue, err := NewQueueConnection(config)
	assert.Nil(t, err)

	message := uuid.Must(uuid.NewV4()).String()
	err = queue.SendMessage([]byte(message))
	assert.Nil(t, err)

	msg, err := queue.ReadMessage()
	require.Nil(t, err)
	require.Equal(t, message, string(msg))
}

// Test de error al leer mensaje
func TestReadMessage_Error(t *testing.T) {
	config := QueueConfig{Host: "localhost", Port: 5672}
	queue, err := NewQueueConnection(config)
	assert.Nil(t, err)

	message := uuid.Must(uuid.NewV4()).String()
	err = queue.SendMessage([]byte(message))
	assert.Nil(t, err)

	defer func() {
		queue, err = NewQueueConnection(config)
		assert.Nil(t, err)

		_, err = queue.ReadMessage()
		assert.Nil(t, err)

		err = queue.Close()
		assert.Nil(t, err)
	}()

	err = queue.Close()
	assert.Nil(t, err)

	msg, err := queue.ReadMessage()
	require.NotNil(t, err)
	require.Nil(t, msg)
}
