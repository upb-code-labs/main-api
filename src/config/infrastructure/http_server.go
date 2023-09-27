package infrastructure

import (
	accounts "github.com/UPB-Code-Labs/main-api/src/accounts/infrastructure"
	courses "github.com/UPB-Code-Labs/main-api/src/courses/infrastructure/http"
	session "github.com/UPB-Code-Labs/main-api/src/session/infrastructure"
	shared "github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var routesGroups = []func(*gin.RouterGroup){
	accounts.StartAccountsRoutes,
	session.StartSessionRoutes,
	courses.StartCoursesRoutes,
}

func StartHTTPServer() {
	engine := gin.Default()
	engine.Use(shared.ErrorHandlerMiddleware())

	isInProductionEnvironment := shared.GetEnvironment().Environment == "production"
	if isInProductionEnvironment {
		gin.SetMode(gin.ReleaseMode)
	}

	// Configure CORS rules
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowOrigins = []string{shared.GetEnvironment().WebClientUrl}
	corsConfig.AllowCredentials = true
	engine.Use(cors.New(corsConfig))

	// Register routes
	baseGroup := engine.Group("/api/v1")
	for _, registerRoutesGroup := range routesGroups {
		registerRoutesGroup(baseGroup)
	}

	engine.Run(":8080")
}
