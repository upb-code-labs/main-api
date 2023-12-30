package http

import (
	"github.com/UPB-Code-Labs/main-api/src/languages/application"
	"github.com/UPB-Code-Labs/main-api/src/languages/infrastructure/implementations"
	sharedInfrastructure "github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
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
		"",
		sharedInfrastructure.WithAuthenticationMiddleware(),
		sharedInfrastructure.WithAuthorizationMiddleware([]string{"teacher", "student"}),
		controllers.HandleGetLanguages,
	)
	langGroup.GET(
		"/:language_uuid/template",
		sharedInfrastructure.WithAuthenticationMiddleware(),
		sharedInfrastructure.WithAuthorizationMiddleware([]string{"teacher", "student"}),
		controllers.HandleDownloadLanguageTemplate,
	)
}
