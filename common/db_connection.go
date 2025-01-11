package common

import (
	"database/sql"
	"log"
)

func ConnectToDatabase() *sql.DB {
	DB_URL := GetDatabaseURL()

	connection, err := sql.Open("postgres", DB_URL)
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	return connection
}
