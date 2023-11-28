package http

import (
	"net/http"

	"github.com/UPB-Code-Labs/main-api/src/laboratories/domain/dtos"

	"github.com/UPB-Code-Labs/main-api/src/laboratories/application"
	"github.com/UPB-Code-Labs/main-api/src/laboratories/infrastructure/requests"
	"github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
	"github.com/gin-gonic/gin"
)

type LaboratoriesController struct {
	UseCases *application.LaboratoriesUseCases
}

func (controller *LaboratoriesController) HandleCreateLaboratory(c *gin.Context) {
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

func (controller *LaboratoriesController) HandleGetLaboratory(c *gin.Context) {
	userUUID := c.GetString("session_uuid")
	laboratoryUUID := c.Param("laboratory_uuid")

	// Validate the laboratory UUID
	if err := infrastructure.GetValidator().Var(laboratoryUUID, "uuid4"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Laboratory UUID is not valid",
		})
		return
	}

	// Get laboratory
	dto := dtos.GetLaboratoryDTO{
		LaboratoryUUID: laboratoryUUID,
		UserUUID:       userUUID,
	}

	laboratory, err := controller.UseCases.GetLaboratory(&dto)
	if err != nil {
		c.Error(err)
		return
	}

	// Return laboratory
	c.JSON(http.StatusOK, laboratory)
}

func (controller *LaboratoriesController) HandleUpdateLaboratory(c *gin.Context) {
	teacherUUID := c.GetString("session_uuid")
	laboratoryUUID := c.Param("laboratory_uuid")

	// Validate the laboratory UUID
	if err := infrastructure.GetValidator().Var(laboratoryUUID, "uuid4"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Laboratory UUID is not valid",
		})
		return
	}

	// Parse request body
	var request requests.UpdateLaboratoryRequest
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

	// Update laboratory
	dto := request.ToDTO(laboratoryUUID, teacherUUID)
	err := controller.UseCases.UpdateLaboratory(dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
