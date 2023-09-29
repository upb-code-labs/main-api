package infrastructure

import (
	"net/http"

	"github.com/UPB-Code-Labs/main-api/src/courses/application"
	"github.com/UPB-Code-Labs/main-api/src/courses/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/courses/infrastructure/requests"
	"github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
	"github.com/gin-gonic/gin"
)

type CoursesController struct {
	UseCases *application.CoursesUseCases
}

func (controller *CoursesController) HandleCreateCourse(c *gin.Context) {
	teacherUUID := c.GetString("session_uuid")

	// Parse request body
	var request requests.CreateCourseRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	// Validate request body
	if err := infrastructure.GetValidator().Struct(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Validation error",
			"errors":  err.Error(),
		})
		return
	}

	// Get a random color for the course
	color, err := controller.UseCases.GetRandomColor()
	if err != nil {
		c.Error(err)
		return
	}

	// Create course
	dto := &dtos.CreateCourseDTO{
		Name:        request.Name,
		TeacherUUID: teacherUUID,
		Color:       *color,
	}

	course, err := controller.UseCases.SaveCourse(dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"uuid":  course.UUID,
		"name":  course.Name,
		"color": course.Color,
	})
}

func (controller *CoursesController) HandleGetInvitationCode(c *gin.Context) {
	teacherUUID := c.GetString("session_uuid")

	// Validate course uuid
	courseUUID := c.Param("course_uuid")
	if err := infrastructure.GetValidator().Var(courseUUID, "uuid4"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid course uuid",
		})
		return
	}

	invitationCode, err := controller.UseCases.GetInvitationCode(dtos.GetInvitationCodeDTO{
		CourseUUID:  courseUUID,
		TeacherUUID: teacherUUID,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": invitationCode,
	})
}
