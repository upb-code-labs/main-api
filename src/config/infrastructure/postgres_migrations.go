package infrastructure

import (
	"log"

	"github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations() {
	driver, err := postgres.WithInstance(infrastructure.GetPostgresConnection(), &postgres.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://sql/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	if err := m.Up(); err != nil {
		// Ignore no changes error
		if err == migrate.ErrNoChange {
			return
		}
		log.Fatal(err.Error())
	}
}
