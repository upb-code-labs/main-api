package http

import (
	"net/http"

	"github.com/UPB-Code-Labs/main-api/src/laboratories/application"
	"github.com/UPB-Code-Labs/main-api/src/laboratories/infrastructure/requests"
	"github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
	"github.com/gin-gonic/gin"
)

type LaboratoriesController struct {
	UseCases *application.LaboratoriesUseCases
}

func (controller *LaboratoriesController) CreateLaboratory(c *gin.Context) {
	teacherUUID := c.GetString("session_uuid")

	// Parse request body
	var request requests.CreateLaboratoryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Request body is not valid",
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

	// Validate due date is after opening date
	openingDate, err1 := infrastructure.ParseISODate(request.OpeningDate)
	dueDate, err2 := infrastructure.ParseISODate(request.DueDate)
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid date format",
		})
		return
	}

	if dueDate.Before(openingDate) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Due date must be after opening date",
		})
		return
	}

	// Create laboratory
	dto := request.ToDTO(teacherUUID)
	laboratory, err := controller.UseCases.CreateLaboratory(dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"uuid": laboratory.UUID,
	})
}