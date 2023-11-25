package http

import (
	courses_implementation "github.com/UPB-Code-Labs/main-api/src/courses/infrastructure/implementations"
	"github.com/UPB-Code-Labs/main-api/src/laboratories/application"
	"github.com/UPB-Code-Labs/main-api/src/laboratories/infrastructure/implementations"
	"github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
	"github.com/gin-gonic/gin"
)

func StartLaboratoriesRoutes(g *gin.RouterGroup) {
	laboratoriesGroup := g.Group("/laboratories")

	useCases := application.LaboratoriesUseCases{
		LaboratoriesRepository: implementations.GetLaboratoriesPostgresRepositoryInstance(),
		CoursesRepository:      courses_implementation.GetCoursesPgRepository(),
	}

	controller := LaboratoriesController{
		UseCases: &useCases,
	}

	laboratoriesGroup.POST("", infrastructure.WithAuthenticationMiddleware(), infrastructure.WithAuthorizationMiddleware([]string{"teacher"}), controller.CreateLaboratory)
}
