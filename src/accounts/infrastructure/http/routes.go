package infrastructure

import (
	"github.com/UPB-Code-Labs/main-api/src/accounts/application"
	"github.com/UPB-Code-Labs/main-api/src/accounts/infrastructure/implementations"
	sharedInfrastructure "github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
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
		sharedInfrastructure.WithAuthenticationMiddleware(),
		sharedInfrastructure.WithAuthorizationMiddleware([]string{"teacher"}),
		controller.HandleSearchStudents,
	)
	accountsGroup.POST(
		"/admins",
		sharedInfrastructure.WithAuthenticationMiddleware(),
		sharedInfrastructure.WithAuthorizationMiddleware([]string{"admin"}),
		controller.HandleRegisterAdmin,
	)
	accountsGroup.GET(
		"/admins",
		sharedInfrastructure.WithAuthenticationMiddleware(),
		sharedInfrastructure.WithAuthorizationMiddleware([]string{"admin"}),
		controller.HandleGetAdmins,
	)
	accountsGroup.POST("/teachers",
		sharedInfrastructure.WithAuthenticationMiddleware(),
		sharedInfrastructure.WithAuthorizationMiddleware([]string{"admin"}),
		controller.HandleRegisterTeacher,
	)
	accountsGroup.PATCH(
		"/password",
		sharedInfrastructure.WithAuthenticationMiddleware(),
		controller.HandleUpdatePassword,
	)
	accountsGroup.PUT(
		"",
		sharedInfrastructure.WithAuthenticationMiddleware(),
		controller.HandleUpdateProfile,
	)
}
