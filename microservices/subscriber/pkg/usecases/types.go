package usecases

const (
	MessageTag = "message"
	ErrorTag   = "error"
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
