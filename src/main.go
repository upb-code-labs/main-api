package main

import (
	config "github.com/UPB-Code-Labs/main-api/src/config/infrastructure"
	shared "github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
)

func main() {
	// Parse environment variables
	shared.GetEnvironment()

	// Connect to database and run migrations
	shared.GetPostgresConnection()
	defer shared.ClosePostgresConnection()
	config.RunMigrations()

	// Start HTTP server
	config.StartHTTPServer()
}
