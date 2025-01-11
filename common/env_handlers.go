package common

import (
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
	DB_URL := os.Getenv("DB_URL")
	if DB_URL == "" {
		log.Fatal("Couldnt get port from env")
	}

	return DB_URL
}

func GetExternalApiURL() string {
	EXTERNAL_API_URL := os.Getenv("EXTERNAL_API_URL")
	if EXTERNAL_API_URL == "" {
		log.Fatal("Couldnt get port from env")
	}

	return EXTERNAL_API_URL
}
