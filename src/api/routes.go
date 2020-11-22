package api

import (
	"wilhelmiina/schedule"

	"github.com/gin-gonic/gin"
)

type errRes struct {
	Message string
	Success bool
}

func getSubjectsRoute(c *gin.Context) {
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
