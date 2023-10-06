package infrastructure

import (
	"github.com/UPB-Code-Labs/main-api/src/courses/application"
	"github.com/UPB-Code-Labs/main-api/src/courses/infrastructure/implementations"
	"github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
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

	coursesGroup.POST(
		"",
		infrastructure.WithAuthenticationMiddleware(),
		infrastructure.WithAuthorizationMiddleware([]string{"teacher"}),
		controller.HandleCreateCourse,
	)

	coursesGroup.GET(
		"",
		infrastructure.WithAuthenticationMiddleware(),
		infrastructure.WithAuthorizationMiddleware([]string{"teacher", "student"}),
		controller.HandleGetEnrolledCourses,
	)

	coursesGroup.GET(
		":course_uuid/invitation-code",
		infrastructure.WithAuthenticationMiddleware(),
		infrastructure.WithAuthorizationMiddleware([]string{"teacher"}),
		controller.HandleGetInvitationCode,
	)

	coursesGroup.POST(
		"/join/:invitation-code",
		infrastructure.WithAuthenticationMiddleware(),
		infrastructure.WithAuthorizationMiddleware([]string{"student"}),
		controller.HandleJoinCourse,
	)

	coursesGroup.PATCH(
		":course_uuid/visibility",
		infrastructure.WithAuthenticationMiddleware(),
		infrastructure.WithAuthorizationMiddleware([]string{"teacher", "student"}),
		controller.HandleChangeCourseVisibility,
	)

	coursesGroup.PATCH(
		":course_uuid/name",
		infrastructure.WithAuthenticationMiddleware(),
		infrastructure.WithAuthorizationMiddleware([]string{"teacher"}),
		controller.HandleChangeCourseName,
	)

	coursesGroup.POST(
		":course_uuid/students",
		infrastructure.WithAuthenticationMiddleware(),
		infrastructure.WithAuthorizationMiddleware([]string{"teacher"}),
		controller.HandleAddStudentToCourse,
	)

	coursesGroup.GET(
		":course_uuid/students",
		infrastructure.WithAuthenticationMiddleware(),
		infrastructure.WithAuthorizationMiddleware([]string{"teacher"}),
		controller.HandleGetStudentsEnrolledInCourse,
	)
}
