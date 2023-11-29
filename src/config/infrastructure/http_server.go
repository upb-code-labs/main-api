package infrastructure

import (
	accounts_http "github.com/UPB-Code-Labs/main-api/src/accounts/infrastructure/http"
	blocks_http "github.com/UPB-Code-Labs/main-api/src/blocks/infrastructure/http"
	courses_http "github.com/UPB-Code-Labs/main-api/src/courses/infrastructure/http"
	laboratories_http "github.com/UPB-Code-Labs/main-api/src/laboratories/infrastructure/http"
	rubrics_http "github.com/UPB-Code-Labs/main-api/src/rubrics/infrastructure/http"
	session_http "github.com/UPB-Code-Labs/main-api/src/session/infrastructure/http"
	shared_infra "github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var routesGroups = []func(*gin.RouterGroup){
	accounts_http.StartAccountsRoutes,
	blocks_http.StartBlocksRoutes,
	session_http.StartSessionRoutes,
	courses_http.StartCoursesRoutes,
	rubrics_http.StartRubricsRoutes,
	laboratories_http.StartLaboratoriesRoutes,
}

func InstanceHttpServer() (r *gin.Engine) {
	engine := gin.Default()
	engine.Use(shared_infra.ErrorHandlerMiddleware())

	isInProductionEnvironment := shared_infra.GetEnvironment().Environment == "production"
	if isInProductionEnvironment {
		gin.SetMode(gin.ReleaseMode)
	}

	// Configure CORS rules
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowOrigins = []string{shared_infra.GetEnvironment().WebClientUrl}
	corsConfig.AllowCredentials = true
	engine.Use(cors.New(corsConfig))

	// Register routes
	baseGroup := engine.Group("/api/v1")
	for _, registerRoutesGroup := range routesGroups {
		registerRoutesGroup(baseGroup)
	}

	return engine
}
