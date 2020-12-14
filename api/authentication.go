package api

import (
	"encoding/json"
	"errors"

	"github.com/Aapeli123/wilhelmiina/user"

	"github.com/gin-gonic/gin"
)

// ErrSessNotFound is thrown when session is not found while looking for session
var ErrSessNotFound = errors.New("Session not found")

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
		c.AbortWithStatusJSON(200, errRes{
			Message: err.Error(),
			Success: false,
		})
		return
	}
	foundUser, err := user.GetUserByName(req.Username)
	if err != nil {
		c.AbortWithStatusJSON(200, errRes{
			Message: err.Error(),
			Success: false,
		})
		return
	}
	succ, err := foundUser.CheckPassword(req.Password)
	if err != nil {
		c.AbortWithStatusJSON(200, errRes{
			Message: err.Error(),
			Success: false,
		})
		return
	}
	if !succ {
		c.AbortWithStatusJSON(200, errRes{
			Message: "Wrong password",
			Success: false,
		})
		return
	}
	sess, err := addSession(foundUser)
	if err != nil {
		c.AbortWithStatusJSON(200, errRes{
			Message: err.Error(),
			Success: false,
		})
		return
	}
	c.SetCookie("SID", sess.SessionID, 0, "/", "", true, true)
	c.JSON(200, authRes{
		Success:   true,
		SessionID: sess.SessionID,
	})
}

type signupReq struct {
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

func isAdminHandler(c *gin.Context) {
	sid, err := c.Cookie("SID")
	if err != nil {
		c.AbortWithStatusJSON(200, errRes{
			Message: err.Error(),
			Success: false,
		})
		return
	}
	// Validate creator session and permissions
	sess, err := getSession(sid)
	if err != nil {
		c.AbortWithStatusJSON(200, errRes{
			Message: err.Error(),
			Success: false,
		})
		return
	}
	u, err := user.GetUser(sess.UserID)
	if err != nil {
		c.AbortWithStatusJSON(200, errRes{
			Message: err.Error(),
			Success: false,
		})
		return
	}
	c.JSON(200, response{
		Success: true,
		Data:    (u.PermissionLevel > 2),
	})
}

func signupHandler(c *gin.Context) {
	var req signupReq
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		c.AbortWithStatusJSON(200, errRes{
			Message: err.Error(),
			Success: false,
		})
		return
	}
	sid, err := c.Cookie("SID")
	if err != nil {
		c.AbortWithStatusJSON(200, errRes{
			Message: err.Error(),
			Success: false,
		})
		return
	}
	// Validate creator session and permissions
	sess, err := getSession(sid)
	if err != nil {
		c.AbortWithStatusJSON(200, errRes{
			Message: err.Error(),
			Success: false,
		})
		return
	}
	userCreator, err := user.GetUser(sess.UserID)
	if err != nil {
		c.AbortWithStatusJSON(200, errRes{
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
	if userCreator.TemporaryAdmin { // Deletes the user
		removeSess(sess.SessionID)
		user.DeleteUser(userCreator.UUID)
	}
	c.JSON(200, signupRes{
		Success: true,
		UUID:    newUser.UUID,
	})
}

type logoutReq struct {
	SessionID string
}

type logoutRes struct {
	Success bool
}

func logoutHandler(c *gin.Context) {
	sid, err := c.Cookie("SID")
	if err != nil {
		c.AbortWithStatusJSON(200, errRes{
			Message: err.Error(),
			Success: false,
		})
		return
	}
	removeSess(sid)
	c.SetCookie("SID", "", 0, "/", "", true, true)
	c.JSON(200, logoutRes{
		Success: true,
	})
}
