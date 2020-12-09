package api

import (
	"encoding/json"
	"net/http"
	"wilhelmiina/schedule"
	"wilhelmiina/user"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// StartServer starts the http server which will serve the api
func StartServer() {
	r := gin.Default()

	r.GET("/ws", websocketHandle)

	r.GET("/subjects", getSubjectsHandler)
	r.GET("/subjects/:id")

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

	startSessionHandler()

	r.Run(":4000")
}

type response struct {
	Success bool
	Data    interface{}
}

func getGroupHandler(c *gin.Context) {
	groupID := c.Param("id")
	group, err := schedule.GetGroup(groupID)
	if err != nil {
		c.AbortWithStatusJSON(500, errRes{Success: false, Message: err.Error()})
	}
	c.JSON(200, response{Data: group, Success: true})

}

func getGroupsForCourseHandler(c *gin.Context) {
	courseID := c.Param("id")
	groups, err := schedule.GetGroupsForCourse(courseID)
	if err != nil {
		c.AbortWithStatusJSON(500, errRes{Success: false, Message: err.Error()})
	}
	c.JSON(200, response{Data: groups, Success: true})
}

func seasonsHandler(c *gin.Context) {
	seasons, err := schedule.GetSeasons()
	if err != nil {
		c.AbortWithStatusJSON(500, errRes{Success: false, Message: err.Error()})
		return
	}
	c.JSON(200, response{Data: seasons, Success: true})
}
func getSeasonHandler(c *gin.Context) {
	seasonID := c.Param("season")
	season, err := schedule.GetSeason(seasonID)
	if err != nil {
		c.AbortWithStatusJSON(500, errRes{Success: false, Message: err.Error()})
		return
	}
	c.JSON(200, response{Data: season, Success: true})
}

func getGroupsForSeasonHandler(c *gin.Context) {
	seasonID := c.Param("season")
	groups, err := schedule.GetGroupsInSeason(seasonID)
	if err != nil {
		c.AbortWithStatusJSON(500, errRes{Success: false, Message: err.Error()})
		return
	}
	c.JSON(200, response{Data: groups, Success: true})
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
