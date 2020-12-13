package api

import (
	"encoding/json"
	"net/http"

	"github.com/Aapeli123/wilhelmiina/schedule"
	"github.com/Aapeli123/wilhelmiina/user"
	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// StartServer starts the http server which will serve the api
func StartServer() {
	r := gin.Default()

	// Allow Cors from anywhere:
	r.Use(cors.Default())

	r.GET("/ws", websocketHandle)

	r.GET("/subjects", getSubjectsHandler)
	r.GET("/subjects/:id", getSubjectHandler)

	r.GET("/seasons", seasonsHandler)
	r.GET("/seasons/:id", getSeasonHandler)
	r.POST("/seasons/create")

	r.POST("/schedule", scheduleHandler)
	r.POST("/schedule/:seasonid", getScheduleForSeasonHandler)

	r.GET("/groups/:season", getGroupsForSeasonHandler)
	r.GET("/group/:id", getGroupHandler)

	r.GET("/courses/:id", getCourseHandler)
	r.GET("/courses/:id/groups", getGroupsForCourseHandler)
	r.GET("/courses", coursesHandler)

	r.POST("/auth/login", loginHandler)
	r.POST("/auth/adduser", signupHandler)

	r.POST("/messages/send", sendMessageHandler)
	r.POST("/messages/getThread", getThreadHandler)
	r.POST("/messages/getMessage", getMessageHandler)
	r.POST("/messages/getThreads", getThreadsHandler)

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
	sess, err := getSession(req.SID)
	if err != nil {
		c.AbortWithStatusJSON(404, errRes{
			Message: err.Error(),
			Success: false,
		})
		return
	}
	userSchedule, err := schedule.GetSchedule(req.ScheduleID)
	if err != nil {
		c.AbortWithStatusJSON(404, errRes{
			Message: err.Error(),
			Success: false,
		})
		return
	}
	user, err := user.GetUser(sess.UserID)
	if err != nil {
		c.AbortWithStatusJSON(404, errRes{
			Message: err.Error(),
			Success: false,
		})
		return
	}
	if !(sess.UserID != userSchedule.OwnerID || !(user.PermissionLevel < 2)) {
		c.AbortWithStatusJSON(403, errRes{
			Message: "You don't have permission view others schedules",
			Success: false,
		})
		return
	}
	c.JSON(200, response{Success: true, Data: userSchedule})
}

func getScheduleForSeasonHandler(c *gin.Context) {
	seasonID := c.Param("seasonid")
	type scheduleReq struct {
		SID string
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
	sess, err := getSession(req.SID)
	if err != nil {
		c.AbortWithStatusJSON(404, errRes{
			Message: err.Error(),
			Success: false,
		})
		return
	}
	userSchedule, err := schedule.GetScheduleForUser(sess.UserID, seasonID)
	if err != nil {
		c.AbortWithStatusJSON(404, errRes{
			Message: err.Error(),
			Success: false,
		})
		return
	}
	c.JSON(200, response{Success: true, Data: userSchedule})

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
