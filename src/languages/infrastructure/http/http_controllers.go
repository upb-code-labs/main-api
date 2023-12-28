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
	languages, err := controller.UseCases.GetLanguages()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"languages": languages,
	})
}

func (controller *LanguagesController) HandleDownloadLanguageTemplate(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}
