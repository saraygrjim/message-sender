package usecases

const (
	MessageTag    = "message"
	ErrorTag      = "error"
	SubscriberTag = "subscriber"
)

type MessageType string

const (
	MessageInfo  MessageType = "info"
	MessageError MessageType = "error"
)

type Message struct {
	Type    MessageType
	Content string
}
