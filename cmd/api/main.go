// @title           Tech Challenge — Oficina Mecânica API
// @version         1.0
// @description     API for managing a mechanical workshop: service orders, customers, vehicles, parts and services.
// @host            localhost:8080
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
package main

import (
	"log"

	"github.com/gabrielcamargo/oficina-api/configs"
	"github.com/gabrielcamargo/oficina-api/internal/api"
	"github.com/gabrielcamargo/oficina-api/internal/infrastructure/database"
)

func main() {
	cfg, err := configs.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if err := database.RunMigrations(cfg.DatabaseURL, "migrations"); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	db, err := database.NewConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	server := api.NewServer(cfg, db)
	if err := server.Run(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
