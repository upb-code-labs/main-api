package http

import (
	"github.com/UPB-Code-Labs/main-api/src/blocks/application"
	"github.com/UPB-Code-Labs/main-api/src/blocks/infrastructure/implementations"
	languagesImplementations "github.com/UPB-Code-Labs/main-api/src/languages/infrastructure/implementations"
	sharedInfrastructure "github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
	"github.com/gin-gonic/gin"
)

func StartBlocksRoutes(g *gin.RouterGroup) {
	blocksGroup := g.Group("/blocks")

	useCases := application.BlocksUseCases{
		BlocksRepository:    implementations.GetBlocksPostgresRepositoryInstance(),
		LanguagesRepository: languagesImplementations.GetLanguagesRepositoryInstance(),
	}

	controller := BlocksController{
		UseCases: &useCases,
	}

	blocksGroup.PATCH(
		"/markdown_blocks/:block_uuid/content",
		sharedInfrastructure.WithAuthenticationMiddleware(),
		sharedInfrastructure.WithAuthorizationMiddleware([]string{"teacher"}),
		controller.HandleUpdateMarkdownBlockContent,
	)

	blocksGroup.PUT(
		"/test_blocks/:block_uuid",
		sharedInfrastructure.WithAuthenticationMiddleware(),
		sharedInfrastructure.WithAuthorizationMiddleware([]string{"teacher"}),
		controller.HandleUpdateTestBlock,
	)

	blocksGroup.DELETE(
		"/markdown_blocks/:block_uuid",
		sharedInfrastructure.WithAuthenticationMiddleware(),
		sharedInfrastructure.WithAuthorizationMiddleware([]string{"teacher"}),
		controller.HandleDeleteMarkdownBlock,
	)

	blocksGroup.DELETE(
		"/test_blocks/:block_uuid",
		sharedInfrastructure.WithAuthenticationMiddleware(),
		sharedInfrastructure.WithAuthorizationMiddleware([]string{"teacher"}),
		controller.HandleDeleteTestBlock,
	)
}
