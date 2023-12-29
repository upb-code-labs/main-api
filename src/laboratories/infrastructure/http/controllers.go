package http

import (
	"fmt"
	"net/http"

	"github.com/UPB-Code-Labs/main-api/src/laboratories/domain/dtos"
	"github.com/gabriel-vasile/mimetype"

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

func (controller *LaboratoriesController) HandleCreateMarkdownBlock(c *gin.Context) {
	teacherUUID := c.GetString("session_uuid")
	laboratoryUUID := c.Param("laboratory_uuid")

	// Validate the laboratory UUID
	if err := infrastructure.GetValidator().Var(laboratoryUUID, "uuid4"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Laboratory UUID is not valid",
		})
		return
	}

	// Create block
	dto := dtos.CreateMarkdownBlockDTO{
		LaboratoryUUID: laboratoryUUID,
		TeacherUUID:    teacherUUID,
	}

	blockUUID, err := controller.UseCases.CreateMarkdownBlock(&dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"uuid": blockUUID,
	})
}

func (controller *LaboratoriesController) HandleCreateTestBlock(c *gin.Context) {
	teacherUUID := c.GetString("session_uuid")
	laboratoryUUID := c.Param("laboratory_uuid")

	// Validate the request struct
	languageUUID := c.PostForm("language_uuid")
	name := c.PostForm("block_name")

	req := requests.CreateTestBlockRequest{
		LaboratoryUUID: laboratoryUUID,
		LanguageUUID:   languageUUID,
		Name:           name,
	}

	if err := infrastructure.GetValidator().Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Validation error",
			"errors":  err.Error(),
		})
		return
	}

	// Validate the test archive
	multipartFile, err := c.FormFile("test_archive")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Please, make sure to send the test archive",
		})
		return
	}

	if multipartFile.Size > infrastructure.GetEnvironment().ArchiveMaxSizeKb*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("The test archive must be smaller than %d KB", infrastructure.GetEnvironment().ArchiveMaxSizeKb),
		})
		return
	}

	file, err := multipartFile.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "There was an error while reading the test archive",
		})
		return
	}
	defer file.Close()

	mtype, err := mimetype.DetectReader(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "There was an error while reading the MIME type of the test archive",
		})
		return
	}

	if mtype.String() != "application/zip" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Please, make sure to send a ZIP archive",
		})
		return
	}

	// Create the DTO
	dto := dtos.CreateTestBlockDTO{
		LaboratoryUUID: laboratoryUUID,
		TeacherUUID:    teacherUUID,
		LanguageUUID:   languageUUID,
		Name:           name,
		MultipartFile:  &file,
	}

	// Create the block
	blockUUID, err := controller.UseCases.CreateTestBlock(&dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"uuid": blockUUID,
	})
}
