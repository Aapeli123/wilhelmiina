package api

import (
	"github.com/Aapeli123/wilhelmiina/schedule"

	"github.com/gin-gonic/gin"
)

type errRes struct {
	Message string
	Success bool
}

type request struct {
	ID string
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
		return
	}
	c.JSON(200, response{Data: group, Success: true})

}

func getSubjectHandler(c *gin.Context) {
	subjectID := c.Param("id")
	subject, err := schedule.GetGroup(subjectID)
	if err != nil {
		c.AbortWithStatusJSON(500, errRes{Success: false, Message: err.Error()})
		return
	}
	c.JSON(200, response{Data: subject, Success: true})
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
	courses, err := schedule.LoadCourses()
	if err != nil {
		c.AbortWithStatusJSON(500, errRes{Success: false, Message: err.Error()})
		return
	}
	c.JSON(200, response{Data: courses, Success: true})
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
	id := c.Param("id")
	course, err := schedule.GetCourse(id)
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
	courses, err := schedule.LoadCourses()
	if err != nil {
		c.AbortWithStatusJSON(404, errRes{
			Message: "Error loading courses: " + err.Error(),
			Success: false,
		})
		return
	}
	c.JSON(200, courses)
}
