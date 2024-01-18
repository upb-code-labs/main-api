package http

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

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

	realTimeUpdater := submissionsImplementations.GetSubmissionsRealTimeUpdatesSenderInstance()

	// Create a channel to send real time updates about the submission
	updatesChannel := realTimeUpdater.CreateChannel(currentStatus.SubmissionUUID)
	defer realTimeUpdater.DeleteChannel(currentStatus.SubmissionUUID)

	/*
		 Add the current status to the channel. Please, note that further updates will be sent
		from the Real Time Updates queue manager
	*/
	go realTimeUpdater.SendUpdate(currentStatus)

	// Connection timeout
	timeoutCh := time.After(5 * time.Minute)

	// Send real time updates
	c.Stream(func(_ io.Writer) bool {
		select {
		// A new update was received
		case update, ok := <-*updatesChannel:
			// Check the channel is not closed
			if !ok {
				return false
			}

			// Parse the update to a JSON
			json, _ := json.Marshal(update)

			// Send the update
			c.SSEvent("update", string(json))

			shouldCloseConnection := update.SubmissionStatus == "ready" &&
				sharedInfrastructure.GetEnvironment().ExecEnvironment == "testing"

			if shouldCloseConnection {
				return false
			} else {
				return true
			}

		// The connection timed out
		case <-timeoutCh:
			return false

		// The client closed the connection
		case <-c.Writer.CloseNotify():
			return false
		}
	})
}

// HandleGetSubmissionArchive controller to handle the request to get the archive with the student's
// code for the given submission
func (controller *SubmissionsController) HandleGetSubmissionArchive(c *gin.Context) {
	userUUID := c.GetString("session_uuid")
	userRole := c.GetString("session_role")
	submissionUUID := c.Param("submission_uuid")

	// Validate the submissionUUID
	if err := sharedInfrastructure.GetValidator().Var(submissionUUID, "uuid4"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "The submission UUID is not valid",
		})
		return
	}

	dto := dtos.GetSubmissionArchiveDTO{
		UserUUID:       userUUID,
		UserRole:       userRole,
		SubmissionUUID: submissionUUID,
	}

	archiveBytes, err := controller.UseCases.GetSubmissionArchive(&dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.Data(http.StatusOK, "application/zip", archiveBytes)
}
