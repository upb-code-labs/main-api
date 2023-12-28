package http

import (
	"net/http"

	"github.com/UPB-Code-Labs/main-api/src/languages/application"
	"github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
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
	// Get the language UUID
	languageUUID := c.Param("language_uuid")

	// Validate the language UUID
	if err := infrastructure.GetValidator().Var(languageUUID, "uuid4"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Language UUID is not valid",
		})
		return
	}

	// Get the language template
	template, err := controller.UseCases.GetLanguageTemplate(languageUUID)
	if err != nil {
		c.Error(err)
		return
	}

	// Return the template
	c.Data(http.StatusOK, "application/zip", template)
}
