package http

import (
	blocksImplementation "github.com/UPB-Code-Labs/main-api/src/blocks/infrastructure/implementations"
	coursesImplementation "github.com/UPB-Code-Labs/main-api/src/courses/infrastructure/implementations"
	"github.com/UPB-Code-Labs/main-api/src/laboratories/application"
	"github.com/UPB-Code-Labs/main-api/src/laboratories/infrastructure/implementations"
	languagesImplementation "github.com/UPB-Code-Labs/main-api/src/languages/infrastructure/implementations"
	rubricImplementation "github.com/UPB-Code-Labs/main-api/src/rubrics/infrastructure/implementations"
	"github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
	staticFilesImplementations "github.com/UPB-Code-Labs/main-api/src/static-files/infrastructure/implementations"
	"github.com/gin-gonic/gin"
)

func StartLaboratoriesRoutes(g *gin.RouterGroup) {
	laboratoriesGroup := g.Group("/laboratories")

	useCases := application.LaboratoriesUseCases{
		StaticFilesRepository:  &staticFilesImplementations.StaticFilesMicroserviceImplementation{},
		LaboratoriesRepository: implementations.GetLaboratoriesPostgresRepositoryInstance(),
		CoursesRepository:      coursesImplementation.GetCoursesPgRepository(),
		RubricsRepository:      rubricImplementation.GetRubricsPgRepository(),
		LanguagesRepository:    languagesImplementation.GetLanguagesRepositoryInstance(),
		BlocksRepository:       blocksImplementation.GetBlocksPostgresRepositoryInstance(),
	}

	controller := LaboratoriesController{
		UseCases: &useCases,
	}

	laboratoriesGroup.POST(
		"",
		infrastructure.WithAuthenticationMiddleware(),
		infrastructure.WithAuthorizationMiddleware([]string{"teacher"}),
		controller.HandleCreateLaboratory,
	)

	laboratoriesGroup.GET(
		"/:laboratory_uuid",
		infrastructure.WithAuthenticationMiddleware(),
		infrastructure.WithAuthorizationMiddleware([]string{"teacher", "student"}),
		controller.HandleGetLaboratory,
	)

	laboratoriesGroup.GET(
		"/:laboratory_uuid/information",
		infrastructure.WithAuthenticationMiddleware(),
		infrastructure.WithAuthorizationMiddleware([]string{"teacher", "student"}),
		controller.HandleGetLaboratoryInformation,
	)

	laboratoriesGroup.PUT(
		"/:laboratory_uuid",
		infrastructure.WithAuthenticationMiddleware(),
		infrastructure.WithAuthorizationMiddleware([]string{"teacher"}),
		controller.HandleUpdateLaboratory,
	)

	laboratoriesGroup.GET(
		"/:laboratory_uuid/progress",
		infrastructure.WithAuthenticationMiddleware(),
		infrastructure.WithAuthorizationMiddleware([]string{"teacher"}),
		controller.HandleGetLaboratoryProgress,
	)

	laboratoriesGroup.GET(
		"/:laboratory_uuid/students/:student_uuid/progress",
		infrastructure.WithAuthenticationMiddleware(),
		infrastructure.WithAuthorizationMiddleware([]string{"teacher", "student"}),
		controller.HandleGetProgressOfStudentInLaboratory,
	)

	laboratoriesGroup.POST(
		"/markdown_blocks/:laboratory_uuid",
		infrastructure.WithAuthenticationMiddleware(),
		infrastructure.WithAuthorizationMiddleware([]string{"teacher"}),
		controller.HandleCreateMarkdownBlock,
	)

	laboratoriesGroup.POST(
		"/test_blocks/:laboratory_uuid",
		infrastructure.WithAuthenticationMiddleware(),
		infrastructure.WithAuthorizationMiddleware([]string{"teacher"}),
		controller.HandleCreateTestBlock,
	)
}
