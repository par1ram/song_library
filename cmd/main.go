package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	api "github.com/par1ram/song-library/api"
	"github.com/par1ram/song-library/common"
	"github.com/pressly/goose"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	PORT := common.GetPort()

	dbCon := common.ConnectToDatabase()
	defer dbCon.Close()

	// Можно выбрать уровень логирования
	apiCfg := api.NewApiConfig(dbCon, logrus.DebugLevel)

	// Миграции
	if err := goose.Up(dbCon, "sql/schema"); err != nil {
		log.Fatalf("Error applying migrations: %v", err)
	}

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	fs := http.FileServer(http.Dir("./docs"))
	router.Handle("/docs/*", http.StripPrefix("/docs/", fs))

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/docs/swagger.yaml"),
	))

	router.Post("/songs/filter", apiCfg.GetSongWithFiltersAndPagination)
	router.Get("/songs/{id}/verses", apiCfg.GetSongVersesWithPagination)

	router.Post("/songs/add", apiCfg.InsertSong)
	router.Put("/songs/update", apiCfg.UpdateSong)
	router.Delete("/songs/delete", apiCfg.DeleteSong)
	router.Patch("/songs/{id}", apiCfg.PatchSong)

	server := &http.Server{
		Addr:           ":" + PORT,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	go func() {
		fmt.Println("Server started on port:", PORT)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exiting")
}
