package http

import (
	"net/http"

	"github.com/UPB-Code-Labs/main-api/src/languages/application"
	"github.com/gin-gonic/gin"
)

type LanguagesController struct {
	UseCases *application.LanguageUseCases
}

func (controller *LanguagesController) HandleGetLanguages(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}

func (controller *LanguagesController) HandleDownloadLanguageTemplate(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}
