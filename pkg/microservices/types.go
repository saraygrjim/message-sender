package microservices

const (
	MessageTag    = "message"
	ErrorTag      = "error"
	SubscriberTag = "subscriber"
)

type Message struct {
	Type    string
	Content string
}
