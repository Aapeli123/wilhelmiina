package api

import (
	"encoding/json"
	"time"
	"wilhelmiina/user"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var sessions []Session

// Session represents an user that has logged in.
// To complete any actions that require authentication, you need an valid session id.
// sessions are stored in a slice and expire in 30 minutes after inactivity if not defined otherwise.
type Session struct {
	SessionID string
	WsConn    *websocket.Conn
	UserID    string
	Expires   int64
}

type authReq struct {
	Username string
	Password string
}

type authRes struct {
	Success   bool
	SessionID string
}

func authHandler(c *gin.Context) {
	var req authReq
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		c.AbortWithStatusJSON(404, errRes{
			Message: err.Error(),
			Success: false,
		})
		return
	}
	foundUser, err := user.GetUserByName(req.Username)
	if err != nil {
		c.AbortWithStatusJSON(404, errRes{
			Message: err.Error(),
			Success: false,
		})
		return
	}
	succ, err := foundUser.CheckPassword(req.Password)
	if err != nil {
		c.AbortWithStatusJSON(404, errRes{
			Message: err.Error(),
			Success: false,
		})
		return
	}
	if !succ {
		c.AbortWithStatusJSON(404, errRes{
			Message: "Wrong password",
			Success: false,
		})
		return
	}
	sess := addSession(foundUser)
	c.JSON(200, authRes{
		Success:   true,
		SessionID: sess.SessionID,
	})
}

// Adds a session for user. Returns the session id
func addSession(u user.User) Session {
	exprire := time.Now().Add(30 * time.Minute)
	s := Session{
		SessionID: uuid.New().String(),
		Expires:   exprire.Unix(),
		UserID:    u.UUID,
	}

	sessions = append(sessions, s)
	return s
}

func startSessionHandler() {
	ticker := time.NewTicker(time.Second * 5)
	go sessionHandler(ticker)
}

func removeSess(index int) {
	sessions = append(sessions[:index], sessions[index+1:]...)
}

func sessionHandler(ticker *time.Ticker) {
	for {
		t := <-ticker.C
		for i, sess := range sessions {
			if sess.Expires < t.Unix() {
				removeSess(i)
			}
		}
	}
}
