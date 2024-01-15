package http

import (
	"github.com/UPB-Code-Labs/main-api/src/grades/application"
	"github.com/UPB-Code-Labs/main-api/src/grades/infrastructure/implementations"
	sharedInfrastructure "github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
	"github.com/gin-gonic/gin"
)

func StartGradesRoutes(g *gin.RouterGroup) {
	gradesGroup := g.Group("/grades")

	useCases := application.GradesUseCases{
		GradesRepository: implementations.GetGradesPostgresRepositoryInstance(),
	}

	controller := &GradesController{
		UseCases: &useCases,
	}

	gradesGroup.GET(
		"/laboratories/:laboratoryUUID",
		sharedInfrastructure.WithAuthenticationMiddleware(),
		sharedInfrastructure.WithAuthorizationMiddleware([]string{"teacher"}),
		controller.GetSummarizedGradesInLaboratory,
	)
}
