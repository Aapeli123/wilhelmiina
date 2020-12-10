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

	msg := messages.NewMessage(sender.UUID, req.Message)

	err = thread.SendMessage(msg)
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
}
