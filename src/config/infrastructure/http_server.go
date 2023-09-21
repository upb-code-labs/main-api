package infrastructure

import (
	"github.com/UPB-Code-Labs/main-api/src/accounts/infrastructure"
	"github.com/gin-gonic/gin"
)

var routesGroups = []func(*gin.RouterGroup){
	infrastructure.StartAccountsRoutes,
}

func StartHTTPServer() {
	engine := gin.Default()
	baseGroup := engine.Group("/api/v1")

	for _, registerRoutesGroup := range routesGroups {
		registerRoutesGroup(baseGroup)
	}

	engine.Run(":8080")
}
