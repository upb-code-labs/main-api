package http

import (
	courses_implementation "github.com/UPB-Code-Labs/main-api/src/courses/infrastructure/implementations"
	"github.com/UPB-Code-Labs/main-api/src/laboratories/application"
	"github.com/UPB-Code-Labs/main-api/src/laboratories/infrastructure/implementations"
	rubrics_implementation "github.com/UPB-Code-Labs/main-api/src/rubrics/infrastructure/implementations"
	"github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
	"github.com/gin-gonic/gin"
)

func StartLaboratoriesRoutes(g *gin.RouterGroup) {
	laboratoriesGroup := g.Group("/laboratories")

	useCases := application.LaboratoriesUseCases{
		LaboratoriesRepository: implementations.GetLaboratoriesPostgresRepositoryInstance(),
		CoursesRepository:      courses_implementation.GetCoursesPgRepository(),
		RubricsRepository:      rubrics_implementation.GetRubricsPgRepository(),
	}

	controller := LaboratoriesController{
		UseCases: &useCases,
	}

	laboratoriesGroup.POST("", infrastructure.WithAuthenticationMiddleware(), infrastructure.WithAuthorizationMiddleware([]string{"teacher"}), controller.CreateLaboratory)
	laboratoriesGroup.PUT("/:laboratory_uuid", infrastructure.WithAuthenticationMiddleware(), infrastructure.WithAuthorizationMiddleware([]string{"teacher"}), controller.UpdateLaboratory)
}
