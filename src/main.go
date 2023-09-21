package main

import (
	"github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
)

func main() {
	// Parse environment variables
	infrastructure.GetEnvironment()

	infrastructure.GetPostgresConnection()
	infrastructure.RunMigrations()
	defer infrastructure.ClosePostgresConnection()
}
