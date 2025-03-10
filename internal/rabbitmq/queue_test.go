package rabbitmq_test

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"message-sender/internal/rabbitmq"
	"os"
	"testing"
)

var testObject *TestObject

type TestObject struct {
	port      int
	container testcontainers.Container
}

func (o *TestObject) Init() error {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "rabbitmq:management",
		ExposedPorts: []string{"5672", "15672"},
		WaitingFor:   wait.ForListeningPort("5672"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return err
	}

	port, err := container.MappedPort(ctx, "5672")
	fmt.Println(port)
	if err != nil {
		return err
	}

	o.container = container
	o.port = port.Int()

	return nil
}

func (o *TestObject) End() error {
	return o.container.Terminate(context.Background())
}

func TestMain(m *testing.M) {
	testObject = &TestObject{}
	err := testObject.Init()
	if err != nil {
		log.Fatal(err)
	}

	code := m.Run()

	err = testObject.End()
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(code)
}

func Test_DummyTest(t *testing.T) {
	config := rabbitmq.QueueConfig{
		Logger: logrus.New(),
		Port:   &testObject.port,
	}

	queue, err := rabbitmq.NewQueueConnection(config)
	require.NotNil(t, queue)
	require.Nil(t, err)
}

func TestNewQueueConnection(t *testing.T) {
	t.Run("Should not be able to connect a queue without Logger", func(t *testing.T) {
		config := rabbitmq.QueueConfig{}
		queue, err := rabbitmq.NewQueueConnection(config)
		require.Nil(t, queue)
		require.NotNil(t, err)
	})
	t.Run("Should be able to connect a queue without Port", func(t *testing.T) {
		config := rabbitmq.QueueConfig{
			Logger: logrus.New(),
		}

		queue, err := rabbitmq.NewQueueConnection(config)
		require.Nil(t, queue)
		require.NotNil(t, err)
	})
	t.Run("Should be able to connect a queue", func(t *testing.T) {
		config := rabbitmq.QueueConfig{
			Logger: logrus.New(),
			Port:   &testObject.port,
		}

		queue, err := rabbitmq.NewQueueConnection(config)
		require.NotNil(t, queue)
		require.Nil(t, err)
	})
}

func TestSendMessage_Success(t *testing.T) {
	config := rabbitmq.QueueConfig{
		Logger: logrus.New(),
		Port:   &testObject.port,
	}
	queue, err := rabbitmq.NewQueueConnection(config)
	require.NotNil(t, queue)
	require.Nil(t, err)

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
	config := rabbitmq.QueueConfig{
		Logger: logrus.New(),
		Port:   &testObject.port,
	}
	queue, err := rabbitmq.NewQueueConnection(config)
	require.NotNil(t, queue)
	require.Nil(t, err)

	err = queue.Close()
	assert.Nil(t, err)

	message := uuid.Must(uuid.NewV4()).String()
	err = queue.SendMessage([]byte(message))
	require.NotNil(t, err)

}

func TestReadMessage_Success(t *testing.T) {
	config := rabbitmq.QueueConfig{
		Logger: logrus.New(),
		Port:   &testObject.port,
	}
	queue, err := rabbitmq.NewQueueConnection(config)
	require.NotNil(t, queue)
	require.Nil(t, err)

	message := uuid.Must(uuid.NewV4()).String()
	err = queue.SendMessage([]byte(message))
	assert.Nil(t, err)

	msg, err := queue.ReadMessage()
	require.Nil(t, err)
	require.Equal(t, message, string(msg))
}

func TestReadMessage_Error(t *testing.T) {
	config := rabbitmq.QueueConfig{
		Logger: logrus.New(),
		Port:   &testObject.port,
	}
	queue, err := rabbitmq.NewQueueConnection(config)
	require.NotNil(t, queue)
	require.Nil(t, err)

	err = queue.Close()
	assert.Nil(t, err)

	msg, err := queue.ReadMessage()
	require.NotNil(t, err)
	require.Nil(t, msg)

}
