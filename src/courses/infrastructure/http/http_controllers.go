package infrastructure

import (
	"net/http"

	"github.com/UPB-Code-Labs/main-api/src/courses/application"
	"github.com/gin-gonic/gin"
)

type CoursesController struct {
	UseCases *application.CoursesUseCases
}

func (controller *CoursesController) HandleCreateCourse(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "To be implemented",
	})
}
