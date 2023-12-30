package http

import (
	"github.com/UPB-Code-Labs/main-api/src/rubrics/application"
	"github.com/UPB-Code-Labs/main-api/src/rubrics/infrastructure/implementations"
	sharedInfrastructure "github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
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
		sharedInfrastructure.WithAuthenticationMiddleware(),
		sharedInfrastructure.WithAuthorizationMiddleware([]string{"teacher"}),
		controller.HandleCreateRubric,
	)

	rubricsGroup.GET(
		"",
		sharedInfrastructure.WithAuthenticationMiddleware(),
		sharedInfrastructure.WithAuthorizationMiddleware([]string{"teacher"}),
		controller.HandleGetRubricsCreatedByTeacher,
	)

	rubricsGroup.GET(
		"/:rubricUUID",
		sharedInfrastructure.WithAuthenticationMiddleware(),
		sharedInfrastructure.WithAuthorizationMiddleware([]string{"teacher"}),
		controller.HandleGetRubricByUUID,
	)

	rubricsGroup.PATCH(
		"/:rubricUUID/name",
		sharedInfrastructure.WithAuthenticationMiddleware(),
		sharedInfrastructure.WithAuthorizationMiddleware([]string{"teacher"}),
		controller.HandleUpdateRubricName,
	)

	rubricsGroup.POST(
		"/:rubricUUID/objectives",
		sharedInfrastructure.WithAuthenticationMiddleware(),
		sharedInfrastructure.WithAuthorizationMiddleware([]string{"teacher"}),
		controller.HandleAddObjectiveToRubric,
	)

	rubricsGroup.POST(
		"/objectives/:objectiveUUID/criteria",
		sharedInfrastructure.WithAuthenticationMiddleware(),
		sharedInfrastructure.WithAuthorizationMiddleware([]string{"teacher"}),
		controller.HandleAddCriteriaToObjective,
	)

	rubricsGroup.PUT(
		"/objectives/:objectiveUUID",
		sharedInfrastructure.WithAuthenticationMiddleware(),
		sharedInfrastructure.WithAuthorizationMiddleware([]string{"teacher"}),
		controller.HandleUpdateObjective,
	)

	rubricsGroup.DELETE(
		"/objectives/:objectiveUUID",
		sharedInfrastructure.WithAuthenticationMiddleware(),
		sharedInfrastructure.WithAuthorizationMiddleware([]string{"teacher"}),
		controller.HandleDeleteObjective,
	)

	rubricsGroup.PUT(
		"/criteria/:criteriaUUID",
		sharedInfrastructure.WithAuthenticationMiddleware(),
		sharedInfrastructure.WithAuthorizationMiddleware([]string{"teacher"}),
		controller.HandleUpdateCriteria,
	)

	rubricsGroup.DELETE(
		"/criteria/:criteriaUUID",
		sharedInfrastructure.WithAuthenticationMiddleware(),
		sharedInfrastructure.WithAuthorizationMiddleware([]string{"teacher"}),
		controller.HandleDeleteCriteria,
	)
}
