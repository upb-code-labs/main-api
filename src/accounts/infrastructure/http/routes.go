package infrastructure

import (
	"github.com/UPB-Code-Labs/main-api/src/accounts/application"
	"github.com/UPB-Code-Labs/main-api/src/accounts/infrastructure/implementations"
	shared_infrastructure "github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
	"github.com/gin-gonic/gin"
)

func StartAccountsRoutes(g *gin.RouterGroup) {
	accountsGroup := g.Group("/accounts")

	useCases := application.AccountsUseCases{
		AccountsRepository: implementations.GetAccountsPgRepository(),
		PasswordsHasher:    implementations.GetArgon2PasswordsHasher(),
	}

	controller := &AccountsController{
		UseCases: &useCases,
	}

	accountsGroup.POST("/students", controller.HandleRegisterStudent)
	accountsGroup.GET(
		"/students",
		shared_infrastructure.WithAuthenticationMiddleware(),
		shared_infrastructure.WithAuthorizationMiddleware([]string{"teacher"}),
		controller.HandleSearchStudents,
	)
	accountsGroup.POST(
		"/admins",
		shared_infrastructure.WithAuthenticationMiddleware(),
		shared_infrastructure.WithAuthorizationMiddleware([]string{"admin"}),
		controller.HandleRegisterAdmin,
	)
	accountsGroup.GET(
		"/admins",
		shared_infrastructure.WithAuthenticationMiddleware(),
		shared_infrastructure.WithAuthorizationMiddleware([]string{"admin"}),
		controller.HandleGetAdmins,
	)
	accountsGroup.POST("/teachers",
		shared_infrastructure.WithAuthenticationMiddleware(),
		shared_infrastructure.WithAuthorizationMiddleware([]string{"admin"}),
		controller.HandleRegisterTeacher,
	)
}
