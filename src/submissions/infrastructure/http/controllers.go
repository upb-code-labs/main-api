package http

import (
	"net/http"

	"github.com/UPB-Code-Labs/main-api/src/submissions/application"
	"github.com/gin-gonic/gin"
)

type SubmissionsController struct {
	UseCases *application.SubmissionUseCases
}

func (controller *SubmissionsController) HandleReceiveSubmissions(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}

func (controller *SubmissionsController) HandleGetSubmission(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}
