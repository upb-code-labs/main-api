package infrastructure

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Environment struct {
	DbConnectionString string `default:"postgres://postgres:postgres@localhost:5432/codelabs?sslmode=disable"`
}

var environment *Environment

func GetEnvironment() *Environment {
	if environment == nil {
		environment = &Environment{}

		err := envconfig.Process("", environment)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	return environment
}
