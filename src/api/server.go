package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// StartServer starts the http server which will serve the api
func StartServer(port int) {
	r := gin.Default()

	r.GET("/ws", websocketHandle)
	r.GET("/subjects", getSubjectsRoute)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func websocketHandle(c *gin.Context) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	// TODO Do something with connection...
	conn.Close() // Close the connection
}
