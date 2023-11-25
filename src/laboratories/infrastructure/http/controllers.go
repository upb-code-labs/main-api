package http

import (
	"net/http"

	"github.com/UPB-Code-Labs/main-api/src/laboratories/application"
	"github.com/gin-gonic/gin"
)

type LaboratoriesController struct {
	UseCases *application.LaboratoriesUseCases
}

func (controller *LaboratoriesController) CreateLaboratory(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}
