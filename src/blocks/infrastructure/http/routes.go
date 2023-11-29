package http

import (
	"github.com/UPB-Code-Labs/main-api/src/blocks/application"
	"github.com/UPB-Code-Labs/main-api/src/blocks/infrastructure/implementations"
	shared_infrastructure "github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
	"github.com/gin-gonic/gin"
)

func StartBlocksRoutes(g *gin.RouterGroup) {
	blocksGroup := g.Group("/blocks")

	useCases := application.BlocksUseCases{
		BlocksRepository: implementations.GetBlocksPostgresRepositoryInstance(),
	}

	controller := BlocksController{
		UseCases: &useCases,
	}

	blocksGroup.PUT(
		"/markdown_blocks/:block_uuid/content",
		shared_infrastructure.WithAuthenticationMiddleware(),
		shared_infrastructure.WithAuthorizationMiddleware([]string{"teacher"}),
		controller.HandleUpdateMarkdownBlockContent,
	)
}
