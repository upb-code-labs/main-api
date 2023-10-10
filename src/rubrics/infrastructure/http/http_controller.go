package http

import (
	"net/http"

	"github.com/UPB-Code-Labs/main-api/src/rubrics/application"
	"github.com/UPB-Code-Labs/main-api/src/rubrics/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/rubrics/infrastructure/requests"
	shared_infrastructure "github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
	"github.com/gin-gonic/gin"
)

type RubricsController struct {
	UseCases *application.RubricsUseCases
}

func (controller *RubricsController) HandleCreateRubric(c *gin.Context) {
	teacher_uuid := c.GetString("session_uuid")

	// Parse request body
	var request requests.CreateRubricRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	// Validate request body
	if err := shared_infrastructure.GetValidator().Struct(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Validation error",
			"errors":  err.Error(),
		})
		return
	}

	// Create DTO
	dto := dtos.CreateRubricDTO{
		TeacherUUID: teacher_uuid,
		Name:        request.Name,
	}

	// Create the course
	rubric, err := controller.UseCases.CreateRubric(&dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Rubric created",
		"uuid":    rubric.UUID,
		"name":    rubric.Name,
	})
}