package api

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

func sendMessageHandler(c *gin.Context) {
	type sendMessageReq struct {
		SID     string
		Title   string
		Message string
	}
	var req sendMessageReq
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		c.AbortWithStatusJSON(500, errRes{Message: err.Error(), Success: false})
	}
}
