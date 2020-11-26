package api

import (
	"encoding/json"
	"wilhelmiina/schedule"

	"github.com/gin-gonic/gin"
)

type errRes struct {
	Message string
	Success bool
}

type request struct {
	ID string
}

func getSubjectsHandler(c *gin.Context) {
	subjects, err := schedule.LoadSubjects()
	if err != nil {
		c.AbortWithStatusJSON(404, errRes{
			Message: "Error loading subjects: " + err.Error(),
			Success: false,
		})
		return
	}
	c.JSON(200, subjects)
}

func getCourseHandler(c *gin.Context) {
	var req request
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		c.AbortWithStatusJSON(500, errRes{
			Message: "Error : " + err.Error(),
			Success: false,
		})
		return
	}
	course, err := schedule.GetCourse(req.ID)
	if err != nil {
		c.AbortWithStatusJSON(404, errRes{
			Message: "Error : " + err.Error(),
			Success: false,
		})
		return
	}
	c.JSON(200, course)
}

func getCoursesHandler(c *gin.Context) {
	courses, err := schedule.LoadSubjects()
	if err != nil {
		c.AbortWithStatusJSON(404, errRes{
			Message: "Error loading courses: " + err.Error(),
			Success: false,
		})
		return
	}
	c.JSON(200, courses)
}
