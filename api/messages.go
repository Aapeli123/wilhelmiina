package api

import (
	"encoding/json"

	"github.com/Aapeli123/wilhelmiina/messages"
	"github.com/Aapeli123/wilhelmiina/user"
	"github.com/gin-gonic/gin"
)

func sendMessageHandler(c *gin.Context) {
	type sendMessageReq struct {
		SID         string
		ThreadTitle string
		Message     string
		Recievers   []string
		ThreadID    string
	}
	var req sendMessageReq
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		c.AbortWithStatusJSON(500, errRes{Message: err.Error(), Success: false})
		return
	}
	sess, err := getSession(req.SID)

	if err != nil {
		c.AbortWithStatusJSON(500, errRes{Message: err.Error(), Success: false})
		return
	}

	sender, err := user.GetUser(sess.UserID)
	if err != nil {
		c.AbortWithStatusJSON(500, errRes{Message: err.Error(), Success: false})
		return
	}
	var thread messages.Thread
	if req.ThreadID == "" {
		msgThread, err := messages.CreateThread(sender.UUID, req.Recievers, req.ThreadTitle)
		if err != nil {
			c.AbortWithStatusJSON(500, errRes{Message: err.Error(), Success: false})
			return
		}
		thread = msgThread
	} else {
		thread, err = messages.GetThread(req.ThreadID)
		if err != nil {
			c.AbortWithStatusJSON(500, errRes{Message: err.Error(), Success: false})
			return
		}
	}

	msg := messages.NewMessage(sender.UUID, req.Message, thread.ThreadID)

	err = thread.SendMessage(msg)
	// TODO Send notification through websocket
	/* for _, id := range thread.Members {
		u, _ := user.GetUser(id)
		sess, _ := sessForUser(u)
		for _, s := range sess {

		}
	} */
	if err != nil {
		c.AbortWithStatusJSON(500, errRes{Message: err.Error(), Success: false})
		return
	}
	c.JSON(200, response{Success: true})
}

func getMessageHandler(c *gin.Context) {
	type getMsgReq struct {
		MessageID string
		SID       string
	}
	var req getMsgReq
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		c.AbortWithStatusJSON(500, errRes{Message: err.Error(), Success: false})
		return
	}
	session, err := getSession(req.SID)
	if err != nil {
		c.AbortWithStatusJSON(500, errRes{Message: err.Error(), Success: false})
		return
	}

	msg, err := messages.GetMessage(req.MessageID)
	if err != nil {
		c.AbortWithStatusJSON(500, errRes{Message: err.Error(), Success: false})
		return
	}
	msgThread, err := messages.GetThread(msg.ThreadID)
	if err != nil {
		c.AbortWithStatusJSON(500, errRes{Message: err.Error(), Success: false})
		return
	}
	perms := false
	for _, userID := range msgThread.Members {
		if userID == session.UserID {
			perms = true
			break
		}
	}

	if !perms {
		c.AbortWithStatusJSON(403, errRes{Message: "You don't have permission to view that message", Success: false})
		return
	}

	c.JSON(200, response{
		Data:    msg,
		Success: true,
	})
}

func getThreadsHandler(c *gin.Context) {
	type getThreadsReq struct {
		SID string
	}
	var req getThreadsReq
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		c.AbortWithStatusJSON(500, errRes{Message: err.Error(), Success: false})
		return
	}
	sess, err := getSession(req.SID)
	if err != nil {
		c.AbortWithStatusJSON(500, errRes{Message: err.Error(), Success: false})
		return
	}

	threads, err := messages.GetThreadsForUser(sess.UserID)
	if err != nil {
		c.AbortWithStatusJSON(500, errRes{Message: err.Error(), Success: false})
		return
	}
	c.JSON(200, response{
		Success: true,
		Data:    threads,
	})
}

func getThreadHandler(c *gin.Context) {
	type getThreadReq struct {
		SID      string
		ThreadID string
	}
	var req getThreadReq
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		c.AbortWithStatusJSON(500, errRes{Message: err.Error(), Success: false})
		return
	}
	sess, err := getSession(req.SID)
	if err != nil {
		c.AbortWithStatusJSON(500, errRes{Message: err.Error(), Success: false})
		return
	}
	thread, err := messages.GetThread(req.ThreadID)
	if err != nil {
		c.AbortWithStatusJSON(500, errRes{Message: err.Error(), Success: false})
		return
	}
	perms := false
	for _, m := range thread.Members {
		if m == sess.UserID {
			perms = true
			break
		}
	}
	if !perms {
		c.AbortWithStatusJSON(403, errRes{Message: "You don't have permission to access that thread", Success: false})
		return
	}
	c.JSON(200, response{
		Success: true,
		Data:    thread,
	})
}
