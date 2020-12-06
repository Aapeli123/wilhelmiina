package api

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type wsSess struct {
	ownerID   string
	sessionID string
	ID        string
	conn      *websocket.Conn
}

var websocketSessions = []wsSess{}

func addWsSession(conn *websocket.Conn, sessID string) (*wsSess, error) {
	sess, err := getSession(sessID)
	if err != nil {
		return nil, err
	}
	wsS := wsSess{
		sessionID: sess.SessionID,
		ownerID:   sess.UserID,
		conn:      conn,
		ID:        uuid.New().String(),
	}
	websocketSessions = append(websocketSessions, wsS)
	return &wsS, nil
}

func removeWsSession(ID string) {
	for i, sess := range websocketSessions {
		if sess.ID == ID {
			websocketSessions = append(websocketSessions[0:i], websocketSessions[i+1:]...)
			return
		}
	}
}

type websocketMessage struct {
	Message   string
	SessionID string
}

func handleWebsocketConnection(sess *wsSess) {
	var message websocketMessage
	for {
		err := sess.conn.ReadJSON(&message)
		if err != nil {
			sess.conn.Close()
			removeWsSession(sess.ID)
			break
		}
		// TODO Do something with recieved websocket message...
		fmt.Println(message.Message)
	}
}
