package http

import (
	"net/http"

	sharedInfrastructure "github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
	"github.com/UPB-Code-Labs/main-api/src/submissions/application"
	"github.com/UPB-Code-Labs/main-api/src/submissions/domain/dtos"
	"github.com/gin-gonic/gin"
)

type SubmissionsController struct {
	UseCases *application.SubmissionUseCases
}

func (controller *SubmissionsController) HandleReceiveSubmissions(c *gin.Context) {
	studentUUID := c.GetString("session_uuid")
	testBlockUUID := c.Param("test_block_uuid")

	// Validate the testBlockUUID
	if err := sharedInfrastructure.GetValidator().Var(testBlockUUID, "uuid4"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "The test block UUID is not valid",
		})
		return
	}

	// Validate the submission archive
	fileMH, err := c.FormFile("submission_archive")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Please, make sure to send the submission archive",
		})
		return
	}

	err = sharedInfrastructure.ValidateMultipartFileHeader(fileMH)
	if err != nil {
		c.Error(err)
		return
	}

	// Create the dto
	file, err := fileMH.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "There was an error while reading the test archive",
		})
		return
	}

	dto := dtos.CreateSubmissionDTO{
		StudentUUID:       studentUUID,
		TestBlockUUID:     testBlockUUID,
		SubmissionArchive: &file,
	}

	// Create the submission
	submissionUUID, err := controller.UseCases.SaveSubmission(&dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"uuid": submissionUUID,
	})
}

func (controller *SubmissionsController) HandleGetSubmission(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}
