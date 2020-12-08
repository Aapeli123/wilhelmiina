package api

import (
	"encoding/json"
	"net/http"
	"wilhelmiina/schedule"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// StartServer starts the http server which will serve the api
func StartServer() {
	r := gin.Default()

	r.GET("/ws", websocketHandle)

	r.GET("/subjects", getSubjectsHandler)
	r.GET("/subject/:id")

	r.GET("/seasons", seasonsHandler)
	r.GET("/season/:id", getSeasonHandler)

	r.POST("/schedule", scheduleHandler)

	r.GET("/course/:id", getCourseHandler)
	r.GET("/courses/:season", getCoursesForSeasonHandler)
	r.GET("/courses", coursesHandler)

	r.POST("/auth/login", loginHandler)
	r.POST("/auth/adduser", signupHandler)

	startSessionHandler()

	r.Run(":4000")
}

func seasonsHandler(c *gin.Context) {
	seasons, err := schedule.GetSeasons()
	if err != nil {
		c.AbortWithStatusJSON(500, errRes{Success: false, Message: err.Error()})
		return
	}
	type seasonsRes struct {
		Seasons []schedule.Season
		Success bool
	}
	c.JSON(200, seasonsRes{Seasons: seasons, Success: true})
}
func getSeasonHandler(c *gin.Context) {
	// TODO Get one season based on id
}

func getCoursesForSeasonHandler(c *gin.Context) {
	// Get all courses in specific season
}

func coursesHandler(c *gin.Context) {
	// TODO Get all courses
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type websocketRes struct {
	Message string
	Success bool
}

func scheduleHandler(c *gin.Context) {
	type scheduleReq struct {
		SID        string
		ScheduleID string
	}
	var req scheduleReq
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		c.AbortWithStatusJSON(404, errRes{
			Message: err.Error(),
			Success: false,
		})
		return
	}
	// TODO Get schedule for user specified in request.
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
