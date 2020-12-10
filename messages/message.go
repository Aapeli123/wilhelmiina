package messages

import (
	"context"
	"time"

	"github.com/Aapeli123/wilhelmiina/database"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

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
	Date      int64
	Content   string
	MessageID string
	ThreadID  string
}

// Thread represents a thread of messages
type Thread struct {
	ThreadID string
	Messages []string
	Title    string
	Members  []string
}

// CreateThread creates a new thread with specified users and saves it to database
func CreateThread(creator string, recievers []string, title string) (Thread, error) {
	members := []string{creator}
	members = append(members, recievers...)
	thread := Thread{
		ThreadID: uuid.New().String(),
		Messages: []string{},
		Members:  members,
		Title:    title,
	}
	_, err := database.DbClient.Database("test").Collection("threads").InsertOne(context.TODO(), thread)
	if err != nil {
		return Thread{}, err
	}
	return thread, nil
}

// DeleteThread removes a thread from database it also removes all messages associated with it
func DeleteThread(threadID string) {
	filter := bson.M{
		"threadid": threadID,
	}
	var thread Thread

	result := database.DbClient.Database("test").Collection("threads").FindOneAndDelete(context.TODO(), filter)
	result.Decode(&thread)
	for _, msg := range thread.Messages { // Remove messages associated with thread
		DeleteMessage(msg)
	}
}

// SendMessage sends a response message to all members of the thread.
//
// Basically saves the message to database and sets the appends its id to the threads messages slice. Saves changes to database
func (t *Thread) SendMessage(message Message) error {
	_, err := database.DbClient.Database("test").Collection("messages").InsertOne(context.TODO(), message)
	if err != nil {
		return err
	}
	t.Messages = append(t.Messages, message.MessageID)
	filter := bson.M{
		"threadid": t.ThreadID,
	}
	database.DbClient.Database("test").Collection("threads").FindOneAndReplace(context.TODO(), filter, *t)
	return nil
}

// AddMember adds new user to the threads Members slice.
func (t *Thread) AddMember(userID string) {
	t.Members = append(t.Members, userID)
	filter := bson.M{
		"threadid": t.ThreadID,
	}
	database.DbClient.Database("test").Collection("threads").FindOneAndReplace(context.TODO(), filter, *t)
}

// RemoveMember removes an user froms the threads Members slice.
func (t *Thread) RemoveMember(userID string) {
	for i, m := range t.Members {
		if m == userID {
			t.Members = append(t.Members[:i], t.Members[i+1:]...)
			break
		}
	}
	filter := bson.M{
		"threadid": t.ThreadID,
	}
	database.DbClient.Database("test").Collection("threads").FindOneAndReplace(context.TODO(), filter, *t)
}

// DeleteMessage removes a specific message from the thread. Also removes the message from database
func (t *Thread) DeleteMessage(messageID string) {
	for i, m := range t.Messages {
		if m == messageID {
			t.Messages = append(t.Messages[:i], t.Messages[i+1:]...)
			break
		}
	}
	filter := bson.M{
		"threadid": t.ThreadID,
	}
	database.DbClient.Database("test").Collection("threads").FindOneAndReplace(context.TODO(), filter, *t)
	filter = bson.M{
		"messageid": messageID,
	}
	database.DbClient.Database("test").Collection("messages").FindOneAndDelete(context.TODO(), filter)
	if len(t.Messages) == 0 { // Delete empty threads
		filter = bson.M{
			"threadid": t.ThreadID,
		}
		database.DbClient.Database("test").Collection("threads").FindOneAndDelete(context.TODO(), filter)
	}
}

// GetMessages gets all messages in the thread
func (t *Thread) GetMessages() ([]Message, error) {
	messages := []Message{}
	for _, m := range t.Messages {
		msg, err := GetMessage(m)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

// GetThread should get a thread based on id
func GetThread(ID string) (Thread, error) {
	thread := Thread{}
	err := database.DbClient.Database("test").Collection("threads").FindOne(context.TODO(), bson.M{
		"threadid": ID,
	}).Decode(&thread)
	if err != nil {
		return Thread{}, err
	}
	return thread, nil
}

// NewMessage creates a new message that can then be added to a thread
func NewMessage(from string, content string, threadID string) Message {
	msg := Message{
		Sender:    from,
		Content:   content,
		MessageID: uuid.New().String(),
		Date:      time.Now().Unix(),
		ThreadID:  threadID,
	}
	return msg
}

// GetMessage gets message from database based on message id
func GetMessage(messageID string) (Message, error) {
	msg := Message{}
	err := database.DbClient.Database("test").Collection("messages").FindOne(context.TODO(), bson.M{
		"messageid": messageID,
	}).Decode(&msg)
	if err != nil {
		return Message{}, err
	}
	return msg, nil
}

// DeleteMessage deletes the message from database. It does not remove it from any threads it is in. For that use Thread.DeleteMessage()
func DeleteMessage(messageID string) {
	filter := bson.M{
		"messageid": messageID,
	}
	database.DbClient.Database("test").Collection("messages").FindOneAndDelete(context.TODO(), filter)
}

// GetThreadsForUser returns all threas for specified user
func GetThreadsForUser(userID string) ([]Thread, error) {
	cur, err := database.DbClient.Database("test").Collection("threads").Find(context.TODO(), bson.D{{}}) // Get all threads
	if err != nil {
		return nil, err
	}
	var threads []Thread
	for cur.Next(context.TODO()) {
		var thread Thread
		err := cur.Decode(&thread)
		if err != nil {
			return nil, err
		}
		threads = append(threads, thread)
	}
	var usersThreads []Thread
	for _, thread := range threads {
		if doesContainUser(userID, thread) {
			usersThreads = append(usersThreads, thread)
		}
	}
	return usersThreads, nil
}

func doesContainUser(userID string, thread Thread) bool {
	for _, u := range thread.Members {
		if u == userID {
			return true
		}
	}
	return false
}
