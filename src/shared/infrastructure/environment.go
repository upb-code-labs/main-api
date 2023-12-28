package infrastructure

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type EnvironmentSpec struct {
	Environment                    string `split_words:"true" default:"development"`
	DbConnectionString             string `split_words:"true" default:"postgres://postgres:postgres@localhost:5432/codelabs?sslmode=disable"`
	DbMigrationsPath               string `split_words:"true" default:"file://sql/migrations"`
	JwtSecret                      string `split_words:"true" default:"default"`
	JwtExpirationHours             int    `split_words:"true" default:"6"`
	WebClientUrl                   string `split_words:"true" default:"http://localhost:5173"`
	StaticFilesMicroserviceAddress string `split_words:"true" default:"http://localhost:8081"`
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
