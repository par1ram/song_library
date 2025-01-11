package common

import (
	"fmt"
	"log"
	"os"
)

func GetPort() string {
	PORT := os.Getenv("PORT")
	if PORT == "" {
		log.Fatal("Couldnt get port from env")
	}

	return PORT
}

func GetDatabaseURL() string {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	if user == "" || host == "" || port == "" || name == "" {
		log.Fatal("Database configuration is incomplete")
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, name, sslmode)
}

func GetExternalApiURL() string {
	EXTERNAL_API_URL := os.Getenv("EXTERNAL_API_URL")
	if EXTERNAL_API_URL == "" {
		log.Fatal("Couldnt get port from env")
	}

	return EXTERNAL_API_URL
}
