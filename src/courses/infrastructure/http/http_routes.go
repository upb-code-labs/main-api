package infrastructure

import (
	"github.com/UPB-Code-Labs/main-api/src/courses/application"
	"github.com/UPB-Code-Labs/main-api/src/courses/infrastructure/implementations"
	"github.com/gin-gonic/gin"
)

func StartCoursesRoutes(g *gin.RouterGroup) {
	coursesGroup := g.Group("/courses")

	useCases := application.CoursesUseCases{
		Repository:              implementations.GetCoursesPgRepository(),
		InvitationCodeGenerator: implementations.GetNanoIdInvitationCodeGenerator(),
	}

	controller := CoursesController{
		UseCases: &useCases,
	}

	coursesGroup.POST("/", controller.HandleCreateCourse)
}
