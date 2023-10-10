package http

import (
	"github.com/UPB-Code-Labs/main-api/src/rubrics/application"
	"github.com/UPB-Code-Labs/main-api/src/rubrics/infrastructure/implementations"
	shared_infrastructure "github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
	"github.com/gin-gonic/gin"
)

func StartRubricsRoutes(g *gin.RouterGroup) {
	rubricsGroup := g.Group("/rubrics")

	useCases := application.RubricsUseCases{
		RubricsRepository: implementations.GetRubricsPgRepository(),
	}

	controller := RubricsController{
		UseCases: &useCases,
	}

	rubricsGroup.POST(
		"",
		shared_infrastructure.WithAuthenticationMiddleware(),
		shared_infrastructure.WithAuthorizationMiddleware([]string{"teacher"}),
		controller.HandleCreateRubric,
	)

	rubricsGroup.GET(
		"",
		shared_infrastructure.WithAuthenticationMiddleware(),
		shared_infrastructure.WithAuthorizationMiddleware([]string{"teacher"}),
		controller.HandleGetRubricsCreatedByTeacher,
	)
}
