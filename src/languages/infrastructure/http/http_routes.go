package http

import (
	"github.com/UPB-Code-Labs/main-api/src/languages/application"
	"github.com/UPB-Code-Labs/main-api/src/languages/infrastructure/implementations"
	shared_infrastructure "github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
	"github.com/gin-gonic/gin"
)

func StartLanguagesRoutes(g *gin.RouterGroup) {
	langGroup := g.Group("/languages")

	useCases := application.LanguageUseCases{
		LanguageRepository: implementations.GetLanguagesRepositoryInstance(),
	}

	controllers := LanguagesController{
		UseCases: &useCases,
	}

	langGroup.GET(
		"/",
		shared_infrastructure.WithAuthenticationMiddleware(),
		shared_infrastructure.WithAuthorizationMiddleware([]string{"teacher", "student"}),
		controllers.HandleGetLanguages,
	)
	langGroup.GET(
		"/:language_uuid/template",
		shared_infrastructure.WithAuthenticationMiddleware(),
		shared_infrastructure.WithAuthorizationMiddleware([]string{"teacher", "student"}),
		controllers.HandleDownloadLanguageTemplate,
	)
}
