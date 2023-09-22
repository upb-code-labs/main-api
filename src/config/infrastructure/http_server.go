package infrastructure

import (
	accounts "github.com/UPB-Code-Labs/main-api/src/accounts/infrastructure"
	shared "github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
	"github.com/gin-gonic/gin"
)

var routesGroups = []func(*gin.RouterGroup){
	accounts.StartAccountsRoutes,
}

func StartHTTPServer() {
	engine := gin.Default()
	engine.Use(shared.ErrorHandlerMiddleware())
	baseGroup := engine.Group("/api/v1")

	for _, registerRoutesGroup := range routesGroups {
		registerRoutesGroup(baseGroup)
	}

	engine.Run(":8080")
}