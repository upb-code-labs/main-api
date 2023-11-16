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

	rubricsGroup.GET(
		"/:rubricUUID",
		shared_infrastructure.WithAuthenticationMiddleware(),
		shared_infrastructure.WithAuthorizationMiddleware([]string{"teacher"}),
		controller.HandleGetRubricByUUID,
	)

	rubricsGroup.PATCH(
		"/:rubricUUID/name",
		shared_infrastructure.WithAuthenticationMiddleware(),
		shared_infrastructure.WithAuthorizationMiddleware([]string{"teacher"}),
		controller.HandleUpdateRubricName,
	)

	rubricsGroup.POST(
		"/:rubricUUID/objectives",
		shared_infrastructure.WithAuthenticationMiddleware(),
		shared_infrastructure.WithAuthorizationMiddleware([]string{"teacher"}),
		controller.HandleAddObjectiveToRubric,
	)

	rubricsGroup.POST(
		"/objectives/:objectiveUUID/criteria",
		shared_infrastructure.WithAuthenticationMiddleware(),
		shared_infrastructure.WithAuthorizationMiddleware([]string{"teacher"}),
		controller.HandleAddCriteriaToObjective,
	)

	rubricsGroup.PUT(
		"/objectives/:objectiveUUID",
		shared_infrastructure.WithAuthenticationMiddleware(),
		shared_infrastructure.WithAuthorizationMiddleware([]string{"teacher"}),
		controller.HandleUpdateObjective,
	)

	rubricsGroup.DELETE(
		"/objectives/:objectiveUUID",
		shared_infrastructure.WithAuthenticationMiddleware(),
		shared_infrastructure.WithAuthorizationMiddleware([]string{"teacher"}),
		controller.HandleDeleteObjective,
	)

	rubricsGroup.PUT(
		"/criteria/:criteriaUUID",
		shared_infrastructure.WithAuthenticationMiddleware(),
		shared_infrastructure.WithAuthorizationMiddleware([]string{"teacher"}),
		controller.HandleUpdateCriteria,
	)
}
