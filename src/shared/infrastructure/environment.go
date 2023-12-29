package infrastructure

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type EnvironmentSpec struct {
	// Execution environment
	Environment string `split_words:"true" default:"development"`

	// Connection strings
	DbConnectionString             string `split_words:"true" default:"postgres://postgres:postgres@localhost:5432/codelabs?sslmode=disable"`
	WebClientUrl                   string `split_words:"true" default:"http://localhost:5173"`
	StaticFilesMicroserviceAddress string `split_words:"true" default:"http://localhost:8081"`

	// PgSQL migration files
	DbMigrationsPath string `split_words:"true" default:"file://sql/migrations"`

	// JWT parameters
	JwtSecret          string `split_words:"true" default:"default"`
	JwtExpirationHours int    `split_words:"true" default:"6"`

	// Configuration parameters
	ArchiveMaxSizeKb int64 `split_words:"true" default:"1024"`
}

var environment *EnvironmentSpec

func GetEnvironment() *EnvironmentSpec {
	if environment == nil {
		environment = &EnvironmentSpec{}

		err := envconfig.Process("", environment)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	return environment
}
