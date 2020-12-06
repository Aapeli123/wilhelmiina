package messages

// Message represents a text message sent from an user to another
type Message struct {
	Sender    string
	Recievers []string
	Date      int64
	Content   string
	MessageID string
}

// Thread represents a thread of messages
type Thread struct {
	ThreadID string
	Messages []string
	Members  string
}
