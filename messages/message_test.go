package messages

import (
	"testing"

	"github.com/Aapeli123/wilhelmiina/database"
	"github.com/Aapeli123/wilhelmiina/user"
)

func TestMessageFlow(t *testing.T) {
	database.Init()
	testSender, _ := user.CreateUser("Testaaja", 4, "Test McTestface", "teakjhndwa@fkla.akdawk", "password1")
	testReciever, _ := user.CreateUser("Testaaja2", 4, "Test McTestface2", "tadwwsaeakjhndwa@fkla.akdawk", "password2")
	thread, err := CreateThread(testSender.UUID, []string{testReciever.UUID}, "Test thread")
	if err != nil {
		t.Fatal(err)
	}
	if thread.Members[0] != testSender.UUID {
		t.Fatal("Threads first member should be the user that started it")
	}
	if len(thread.Members) != 2 {
		t.Fatal("There should be 2 users in thread but there were:", len(thread.Members))
	}
	thread2, err := GetThread(thread.ThreadID)
	if err != nil {
		t.Fatal(err)
	}
	if thread.ThreadID != thread2.ThreadID {
		t.Fatal("Returned wrong thread")
	}

	msg := NewMessage(testSender.UUID, "Lorem ipsum etc..")
	err = thread.SendMessage(msg)
	if err != nil {
		t.Fatal(err)
	}

	msg2 := NewMessage(testReciever.UUID, "Also testing")
	thread.SendMessage(msg2)
	thread3, err := GetThread(thread.ThreadID)
	if thread3.Messages[0] != msg.MessageID {
		t.Fatal("First message was not expected")
	}

	gottenMessage, err := GetMessage(thread3.Messages[1])
	if err != nil {
		t.Fatal(err)
	}

	if gottenMessage.Content != "Also testing" {
		t.Fatal("Got wrong message")
	}
	thread3.DeleteMessage(gottenMessage.MessageID)
	msgs, err := thread3.GetMessages()
	if err != nil {
		t.Fatal(err)
	}

	if len(msgs) != 1 {
		t.Fatal("There was an unexpected amount of messages")
	}
	DeleteThread(thread.ThreadID)
	_, err = GetThread(thread.ThreadID)
	if err == nil {
		t.Fatal("Did not error while getting non existant thread")
	}

	_, err = GetMessage(msgs[0].MessageID)
	if err == nil {
		t.Fatal("DeleteThread did not delete messages from thread")
	}
	user.DeleteUser(testSender.UUID)
	user.DeleteUser(testReciever.UUID)
	database.Close()
}
