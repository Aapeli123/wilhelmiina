package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Aapeli123/wilhelmiina/schedule"
	"github.com/Aapeli123/wilhelmiina/user"
	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// StartServer starts the http server which will serve the api
func StartServer() {
	r := gin.Default()
	c := cors.New(cors.Config{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowAllOrigins:  false,
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		AllowHeaders:     []string{"content-type", "content-length"},
	})
	// Allow Cors from anywhere:
	r.Use(c)

	r.GET("/cookies", func(c *gin.Context) {
		fmt.Println(c.Cookie("SID"))
	})

	r.GET("/ws", websocketHandle)

	r.GET("/subjects", getSubjectsHandler)
	r.GET("/subjects/:id", getSubjectHandler)

	r.GET("/seasons", seasonsHandler)
	r.GET("/seasons/:id", getSeasonHandler)
	r.POST("/seasons/create")

	r.POST("/schedule", scheduleHandler)
	r.GET("/schedule/:seasonid", getScheduleForSeasonHandler)

	r.GET("/groups/:season", getGroupsForSeasonHandler)
	r.GET("/group/:id", getGroupHandler)

	r.GET("/courses/:id", getCourseHandler)
	r.GET("/courses/:id/groups", getGroupsForCourseHandler)
	r.GET("/courses", coursesHandler)

	r.POST("/auth/login", loginHandler)
	r.POST("/auth/adduser", signupHandler)
	r.GET("/auth/logout", logoutHandler)

	r.GET("/user", getUserHandler)
	r.GET("/validateSession", sessionValidityHandler)
	r.POST("/messages/send", sendMessageHandler)
	r.POST("/messages/getThread", getThreadHandler)
	r.POST("/messages/getMessage", getMessageHandler)
	r.POST("/messages/getThreads", getThreadsHandler)

	r.GET("/admin", isAdminHandler)

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

func getUserHandler(c *gin.Context) {
	sid, err := c.Cookie("SID")
	if err != nil {
		c.AbortWithStatusJSON(200, errRes{
			Message: err.Error(),
			Success: false,
		})
		return
	}
	sess, err := getSession(sid)
	if err != nil {
		c.AbortWithStatusJSON(200, errRes{
			Message: err.Error(),
			Success: false,
		})
		return
	}
	user, err := user.GetUser(sess.UserID)
	if err != nil {
		c.AbortWithStatusJSON(200, errRes{
			Message: err.Error(),
			Success: false,
		})
		return
	}
	type userRes struct {
		Username string
		UUID     string
		Fullname string
	}
	c.JSON(200, response{
		Data:    user,
		Success: true,
	})
}

func scheduleHandler(c *gin.Context) {
	type scheduleReq struct {
		ScheduleID string
	}
	var req scheduleReq
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
	sess, err := getSession(sid)
	if err != nil {
		c.AbortWithStatusJSON(200, errRes{
			Message: err.Error(),
			Success: false,
		})
		return
	}
	userSchedule, err := schedule.GetSchedule(req.ScheduleID)
	if err != nil {
		c.AbortWithStatusJSON(200, errRes{
			Message: err.Error(),
			Success: false,
		})
		return
	}
	user, err := user.GetUser(sess.UserID)
	if err != nil {
		c.AbortWithStatusJSON(200, errRes{
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
	sid, err := c.Cookie("SID")
	if err != nil {
		c.AbortWithStatusJSON(200, errRes{
			Message: err.Error(),
			Success: false,
		})
		return
	}
	sess, err := getSession(sid)
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
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		conn.WriteJSON(websocketRes{
			Success: false,
			Message: err.Error(),
		})
		conn.Close()
		return
	}
	sid, err := c.Cookie("SID")
	if err != nil {
		conn.WriteJSON(websocketRes{
			Success: false,
			Message: err.Error(),
		})
		conn.Close()
		return
	}
	wsSession, err := addWsSession(conn, sid)
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
