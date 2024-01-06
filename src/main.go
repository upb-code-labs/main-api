package main

import (
	config "github.com/UPB-Code-Labs/main-api/src/config/infrastructure"
	shared "github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
	submissionsImplementations "github.com/UPB-Code-Labs/main-api/src/submissions/infrastructure/implementations"
)

func main() {
	// Parse environment variables
	shared.GetEnvironment()

	// Connect to database and run migrations
	shared.GetPostgresConnection()
	defer shared.ClosePostgresConnection()
	config.RunMigrations()

	// Connect to RabbitMQ
	shared.ConnectToRabbitMQ()
	defer shared.CloseRabbitMQConnection()

	// Start listening for SSE connections
	realTimeSubmissionsUpdatesSender := submissionsImplementations.GetSubmissionsRealTimeUpdatesSenderInstance()
	go realTimeSubmissionsUpdatesSender.Listen()

	// Start HTTP server
	router := config.InstanceHttpServer()
	router.Run(":8080")
}
