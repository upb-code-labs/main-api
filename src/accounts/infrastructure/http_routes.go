package infrastructure

import (
	"github.com/UPB-Code-Labs/main-api/src/accounts/application"
	"github.com/gin-gonic/gin"
)

func StartAccountsRoutes(g *gin.RouterGroup) {
	accountsGroup := g.Group("/accounts")

	useCases := application.AccountsUseCases{
		AccountsRepository: GetAccountsPgRepository(),
		PasswordsHasher:    GetArgon2PasswordsHasher(),
	}

	controller := &AccountsController{
		UseCases: &useCases,
	}

	accountsGroup.POST("/students", controller.HandleRegisterStudent)
}
