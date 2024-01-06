package http

import (
	blocksImplementations "github.com/UPB-Code-Labs/main-api/src/blocks/infrastructure/implementations"
	sharedInfrastructure "github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
	"github.com/UPB-Code-Labs/main-api/src/submissions/application"
	"github.com/UPB-Code-Labs/main-api/src/submissions/infrastructure/implementations"
	"github.com/gin-gonic/gin"
)

func StartSubmissionsRoutes(g *gin.RouterGroup) {
	submissionsGroup := g.Group("/submissions")

	useCases := application.SubmissionUseCases{
		BlocksRepository:        blocksImplementations.GetBlocksPostgresRepositoryInstance(),
		SubmissionsRepository:   implementations.GetSubmissionsRepositoryInstance(),
		SubmissionsQueueManager: implementations.GetSubmissionsRabbitMQQueueManagerInstance(),
	}

	controllers := SubmissionsController{
		UseCases: &useCases,
	}

	submissionsGroup.POST(
		":test_block_uuid",
		sharedInfrastructure.WithAuthenticationMiddleware(),
		sharedInfrastructure.WithAuthorizationMiddleware([]string{"student"}),
		controllers.HandleReceiveSubmissions,
	)

	submissionsGroup.GET(
		"/:test_block_uuid/status",
		sharedInfrastructure.WithAuthenticationMiddleware(),
		sharedInfrastructure.WithAuthorizationMiddleware([]string{"student"}),
		sharedInfrastructure.WithServerSentEventsMiddleware(),
		controllers.HandleGetSubmission,
	)
}
