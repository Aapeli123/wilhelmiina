package messages

import "wilhelmiina/user"

// Message represents a text message sent from an user to another
type Message struct {
	Sender    user.User
	Recievers []user.User
	Date      int64
	Content   string
}
