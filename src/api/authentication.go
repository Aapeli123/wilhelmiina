package api

import (
	"encoding/json"
	"errors"
	"time"
	"wilhelmiina/user"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var sessions []Session

// ErrSessNotFound is thrown when session is not found while looking for session
var ErrSessNotFound = errors.New("Session not found")

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

func loginHandler(c *gin.Context) {
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

type signupReq struct {
	SID           string
	RealName      string
	Username      string
	Email         string
	Password      string
	PermissionLvl int
}
type signupRes struct {
	Success bool
	UUID    string
}

func signupHandler(c *gin.Context) {
	var req signupReq
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		c.AbortWithStatusJSON(404, errRes{
			Message: err.Error(),
			Success: false,
		})
		return
	}

	// Validate creator session and permissions
	sess, err := getSession(req.SID)
	if err != nil {
		c.AbortWithStatusJSON(404, errRes{
			Message: err.Error(),
			Success: false,
		})
		return
	}
	userCreator, err := user.GetUser(sess.UserID)
	if err != nil {
		c.AbortWithStatusJSON(404, errRes{
			Message: err.Error(),
			Success: false,
		})
		return
	}
	if userCreator.PermissionLevel < 3 {
		c.AbortWithStatusJSON(400, errRes{
			Message: "You don't have permissions to do that",
			Success: false,
		})
		return
	}

	if userCreator.PermissionLevel < req.PermissionLvl {
		c.AbortWithStatusJSON(400, errRes{
			Message: "You can't create users with more permissions than you",
			Success: false,
		})
		return
	}

	// Create the user
	newUser, err := user.CreateUser(req.Username, req.PermissionLvl, req.RealName, req.Email, req.Password)
	if err != nil {
		c.AbortWithStatusJSON(500, errRes{
			Message: err.Error(),
			Success: false,
		})
		return
	}
	c.JSON(200, signupRes{
		Success: true,
		UUID:    newUser.UUID,
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

func getSession(SID string) (Session, error) {
	for _, s := range sessions {
		if s.SessionID == SID {
			return s, nil
		}
	}
	return Session{}, ErrSessNotFound
}
