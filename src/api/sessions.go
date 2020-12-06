package api

import (
	"context"
	"log"
	"time"
	"wilhelmiina/database"
	"wilhelmiina/user"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Session represents an user that has logged in.
// To complete any actions that require authentication, you need an valid session id.
// sessions are stored in the database and expire in 30 minutes after inactivity if not defined otherwise.
type Session struct {
	SessionID string
	UserID    string
	Expires   int64
}

// Adds a session for user. Returns the session id
func addSession(u user.User) (Session, error) {
	exprire := time.Now().Add(30 * time.Minute)
	s := Session{
		SessionID: uuid.New().String(),
		Expires:   exprire.Unix(),
		UserID:    u.UUID,
	}
	collection := database.DbClient.Database("test").Collection("sessions")
	_, err := collection.InsertOne(context.TODO(), s)
	if err != nil {
		return Session{}, err
	}
	return s, nil
}

func startSessionHandler() {
	ticker := time.NewTicker(time.Second * 5)
	go sessionHandler(ticker)
}

func removeSess(SID string) {
	collection := database.DbClient.Database("test").Collection("sessions")
	collection.FindOneAndDelete(context.TODO(), bson.M{"sessionid": SID})
}

func sessionHandler(ticker *time.Ticker) {
	for {
		t := <-ticker.C
		sessions, err := getSessions()
		if err != nil {
			log.Println("Error while getting sessions", err)
		}
		for _, sess := range sessions {
			if sess.Expires < t.Unix() {
				removeSess(sess.SessionID)
			}
		}
	}
}

func getSession(SID string) (Session, error) {
	filter := bson.M{
		"sessionid": SID,
	}
	var sess Session
	err := database.DbClient.Database("test").Collection("sessions").FindOne(context.TODO(), filter).Decode(&sess)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return Session{}, ErrSessNotFound
		}
		return Session{}, err
	}
	return sess, nil
}

func getSessions() ([]Session, error) {
	sessions := []Session{}
	cur, err := database.DbClient.Database("test").Collection("sessions").Find(context.TODO(), bson.D{{}})
	if err != nil {
		return nil, err
	}
	for cur.Next(context.TODO()) {
		var elem Session
		cur.Decode(&elem)
		sessions = append(sessions, elem)
	}
	return sessions, nil
}

// validateSession returns true if SID is valid, false otherwise
func validateSession(SID string) bool {
	sess, err := getSession(SID)

	if err != nil {
		return false
	}
	return sess.Expires > time.Now().Unix()
}
