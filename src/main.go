package main

import (
	"github.com/UPB-Code-Labs/main-api/src/config/infrastructure"
)

func main() {
	// Parse environment variables
	infrastructure.GetEnvironment()

	// Connect to database and run migrations
	infrastructure.GetPostgresConnection()
	defer infrastructure.ClosePostgresConnection()
	infrastructure.RunMigrations()

	// Start HTTP server
	infrastructure.StartHTTPServer()
}
