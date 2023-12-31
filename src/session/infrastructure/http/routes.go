package http

import (
	accounts_impl "github.com/UPB-Code-Labs/main-api/src/accounts/infrastructure/implementations"
	sharedInfrastructure "github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"

	"github.com/UPB-Code-Labs/main-api/src/session/application"
	"github.com/gin-gonic/gin"
)

func StartSessionRoutes(g *gin.RouterGroup) {
	sessionGroup := g.Group("/session")

	useCases := application.SessionUseCases{
		AccountsRepository: accounts_impl.GetAccountsPgRepository(),
		PasswordHasher:     accounts_impl.GetArgon2PasswordsHasher(),
		TokenHandler:       sharedInfrastructure.GetJwtTokenHandler(),
	}

	controllers := &SessionControllers{
		UseCases: &useCases,
	}

	sessionGroup.POST("/login", controllers.HandleLogin)

	sessionGroup.DELETE(
		"/logout",
		sharedInfrastructure.WithAuthenticationMiddleware(),
		controllers.HandleLogout,
	)

	sessionGroup.GET(
		"/whoami",
		sharedInfrastructure.WithAuthenticationMiddleware(),
		controllers.HandleWhoAmI,
	)
}
