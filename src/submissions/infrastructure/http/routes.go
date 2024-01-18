package http

import (
	blocksImplementations "github.com/UPB-Code-Labs/main-api/src/blocks/infrastructure/implementations"
	laboratoriesImplementation "github.com/UPB-Code-Labs/main-api/src/laboratories/infrastructure/implementations"
	sharedInfrastructure "github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
	staticFilesImplementations "github.com/UPB-Code-Labs/main-api/src/static-files/infrastructure/implementations"
	"github.com/UPB-Code-Labs/main-api/src/submissions/application"
	"github.com/UPB-Code-Labs/main-api/src/submissions/infrastructure/implementations"
	"github.com/gin-gonic/gin"
)

func StartSubmissionsRoutes(g *gin.RouterGroup) {
	submissionsGroup := g.Group("/submissions")

	useCases := application.SubmissionUseCases{
		StaticFilesRepository:   &staticFilesImplementations.StaticFilesMicroserviceImplementation{},
		LaboratoriesRepository:  laboratoriesImplementation.GetLaboratoriesPostgresRepositoryInstance(),
		BlocksRepository:        blocksImplementations.GetBlocksPostgresRepositoryInstance(),
		SubmissionsRepository:   implementations.GetSubmissionsRepositoryInstance(),
		SubmissionsQueueManager: implementations.GetSubmissionsRabbitMQQueueManagerInstance(),
	}

	controllers := SubmissionsController{
		UseCases: &useCases,
	}

	submissionsGroup.POST(
		"/test_blocks/:test_block_uuid",
		sharedInfrastructure.WithAuthenticationMiddleware(),
		sharedInfrastructure.WithAuthorizationMiddleware([]string{"student"}),
		controllers.HandleReceiveSubmissions,
	)

	submissionsGroup.GET(
		"/test_blocks/:test_block_uuid/status",
		sharedInfrastructure.WithAuthenticationMiddleware(),
		sharedInfrastructure.WithAuthorizationMiddleware([]string{"student"}),
		sharedInfrastructure.WithServerSentEventsMiddleware(),
		controllers.HandleGetSubmission,
	)

	submissionsGroup.GET(
		"/:submission_uuid/archive",
		sharedInfrastructure.WithAuthenticationMiddleware(),
		sharedInfrastructure.WithAuthorizationMiddleware([]string{"teacher", "student"}),
		controllers.HandleGetSubmissionArchive,
	)
}
