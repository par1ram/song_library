package api

import (
	"database/sql"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/par1ram/song-library/internal/database"
	"github.com/sirupsen/logrus"
)

type ApiConfig struct {
	DB         *database.Queries
	HTTPClient *http.Client
	Logger     *logrus.Logger
}

func NewApiConfig(con *sql.DB, logLevel logrus.Level) *ApiConfig {
	logger := logrus.New()
	logger.SetLevel(logLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})

	return &ApiConfig{
		DB:         database.New(con),
		HTTPClient: &http.Client{},
		Logger:     logger,
	}
}
