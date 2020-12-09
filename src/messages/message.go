package messages

// Message represents a text message sent from an user to another
//
// How to send a message:
//
// 1. Create the thread with all the users that should recieve the message
// 2. Create the message with NewMessage and data from user
// 3. Call the threads SendMessage method with the created message
// 4. Profit.
type Message struct {
	Sender    string
	Recievers []MessageReciever
	Date      int64
	Title     string
	Content   string
	MessageID string
}

// MessageReciever represents a user that will recieve the message.
type MessageReciever struct {
	UUID   string
	ReadBy bool
}

// Thread represents a thread of messages
type Thread struct {
	ThreadID string
	Messages []string
	Members  []string
}

// SendMessage sends a response message to all members of the thread.
//
// Basically saves the message to database and sets the appends its id to the threads messages slice. Saves changes to database
func (t *Thread) SendMessage(message Message) error {
	// TODO
}

// AddMember adds new user to the threads Members slice.
func (t *Thread) AddMember(userID string) error {
	// TODO
}

// DeleteMessage removes a specific message from the thread. Also removes the message from database
func (t *Thread) DeleteMessage(messageID string) {
	// TODO
}

// NewMessage creates a new user that can then be added to a thread
func NewMessage(from string, content string, title string) (Message, error) {
	// TODO
}

// DeleteMessage deletes the message from database. It does not remove it from any threads it is in. For that use Thread.DeleteMessage()
func DeleteMessage(messageID string) {
	// TODO
}
