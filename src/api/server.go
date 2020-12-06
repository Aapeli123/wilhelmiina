package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// StartServer starts the http server which will serve the api
func StartServer() {
	r := gin.Default()

	r.GET("/ws", websocketHandle)
	r.GET("/subjects", getSubjectsHandler)
	r.POST("/course", getCourseHandler)
	r.POST("/auth/login", loginHandler)
	r.POST("/auth/adduser", signupHandler)

	startSessionHandler()

	r.Run(":4000")
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type websocketRes struct {
	Message string
	Success bool
}

func websocketHandle(c *gin.Context) {
	type wsAuthReq struct {
		SessionID string
	}
	req := wsAuthReq{}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	err = conn.ReadJSON(&req)
	if err != nil {
		conn.WriteJSON(websocketRes{
			Success: false,
			Message: err.Error(),
		})
		conn.Close()
		return
	}
	wsSession, err := addWsSession(conn, req.SessionID)
	if err != nil {
		conn.WriteJSON(websocketRes{
			Success: false,
			Message: err.Error(),
		})
		conn.Close()
		return
	}
	handleWebsocketConnection(wsSession)
	// conn.Close() // Close the connection
}
