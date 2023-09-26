package infrastructure

import (
	accounts_infrastructure "github.com/UPB-Code-Labs/main-api/src/accounts/infrastructure"
	shared_infrastructure "github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"

	"github.com/UPB-Code-Labs/main-api/src/session/application"
	"github.com/gin-gonic/gin"
)

func StartSessionRoutes(g *gin.RouterGroup) {
	sessionGroup := g.Group("/session")

	useCases := application.SessionUseCases{
		AccountsRepository: accounts_infrastructure.GetAccountsPgRepository(),
		PasswordHasher:     accounts_infrastructure.GetArgon2PasswordsHasher(),
		TokenHandler:       shared_infrastructure.GetJwtTokenHandler(),
	}

	controllers := &SessionControllers{
		UseCases: &useCases,
	}

	sessionGroup.POST("/login", controllers.HandleLogin)

	sessionGroup.DELETE(
		"/logout",
		shared_infrastructure.WithAuthenticationMiddleware(),
		controllers.HandleLogout,
	)

	sessionGroup.GET(
		"/whoami",
		shared_infrastructure.WithAuthenticationMiddleware(),
		controllers.HandleWhoAmI,
	)
}
