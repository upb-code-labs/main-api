package http

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	sharedInfrastructure "github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
	"github.com/UPB-Code-Labs/main-api/src/submissions/application"
	"github.com/UPB-Code-Labs/main-api/src/submissions/domain/dtos"
	submissionsImplementations "github.com/UPB-Code-Labs/main-api/src/submissions/infrastructure/implementations"
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
	studentUUID := c.GetString("session_uuid")
	testBlockUUID := c.Param("test_block_uuid")

	// Validate the testBlockUUID
	if err := sharedInfrastructure.GetValidator().Var(testBlockUUID, "uuid4"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "The test block UUID is not valid",
		})
		return
	}

	// Get the current status of the submission
	currentStatus, err := controller.UseCases.GetSubmissionStatus(studentUUID, testBlockUUID)
	if err != nil {
		c.Error(err)
		return
	}

	// Create a channel to send real time updates about the submission
	realTimeUpdater := submissionsImplementations.GetSubmissionsRealTimeUpdatesSenderInstance()
	updatesChannel := realTimeUpdater.CreateChannel(currentStatus.SubmissionUUID)
	defer realTimeUpdater.DeleteChannel(currentStatus.SubmissionUUID)

	// Add the current status to the channel. Please, note that further updates will be sent
	// from the Real Time Updates queue manager
	go realTimeUpdater.SendUpdate(currentStatus)

	// Send real time updates
	c.Stream(func(_ io.Writer) bool {
		select {
		// A new update was received
		case update := <-*updatesChannel:
			// Parse the update to a JSON
			json, err := json.Marshal(update)

			if err != nil {
				log.Printf(
					"[SSE]: Error while parsing the update to JSON: %s",
					err.Error(),
				)

				realTimeUpdater.DeleteChannel(testBlockUUID)
				return false
			}

			// Send the update
			c.SSEvent("update", string(json))

			// Check if the submission is finished
			if update.SubmissionStatus == "ready" {
				realTimeUpdater.DeleteChannel(testBlockUUID)
				return false
			}

			return true
		// The client closed the connection
		case <-c.Writer.CloseNotify():
			realTimeUpdater.DeleteChannel(testBlockUUID)
			return false
		}
	})
}
