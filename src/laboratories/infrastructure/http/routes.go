package http

import (
	"github.com/UPB-Code-Labs/main-api/src/laboratories/application"
	"github.com/UPB-Code-Labs/main-api/src/laboratories/infrastructure/implementations"
	"github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
	"github.com/gin-gonic/gin"
)

func StartLaboratoriesRoutes(g *gin.RouterGroup) {
	laboratoriesGroup := g.Group("/laboratories")

	useCases := application.LaboratoriesUseCases{
		Repository: implementations.GetLaboratoriesPostgresRepositoryInstance(),
	}

	controller := LaboratoriesController{
		UseCases: &useCases,
	}

	laboratoriesGroup.POST("", infrastructure.WithAuthenticationMiddleware(), infrastructure.WithAuthorizationMiddleware([]string{"teacher"}), controller.CreateLaboratory)
}
