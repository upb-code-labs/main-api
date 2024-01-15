package http

import (
	"net/http"

	"github.com/UPB-Code-Labs/main-api/src/grades/application"
	"github.com/gin-gonic/gin"
)

// GradesController controller to handle the requests to the `/grades` endpoints
type GradesController struct {
	UseCases *application.GradesUseCases
}

// GetSummarizedGradesInLaboratory controller to get the summarized grades of the students in a laboratory
func (controller *GradesController) GetSummarizedGradesInLaboratory(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}
