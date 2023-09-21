package infrastructure

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var pgConnection *sql.DB

func GetPostgresConnection() *sql.DB {
	if pgConnection == nil {
		// Connect
		pgConnectionString := GetEnvironment().DbConnectionString
		db, err := sql.Open("postgres", pgConnectionString)
		if err != nil {
			log.Fatal(err.Error())
		}

		// Check connection
		if err = db.Ping(); err != nil {
			log.Fatal(err.Error())
		}

		// Set connection
		pgConnection = db
	}

	return pgConnection
}

func ClosePostgresConnection() {
	if pgConnection != nil {
		pgConnection.Close()
	}
}
