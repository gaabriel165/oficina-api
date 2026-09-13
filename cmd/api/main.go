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
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gabrielcamargo/oficina-api/configs"
	"github.com/gabrielcamargo/oficina-api/internal/api"
	"github.com/gabrielcamargo/oficina-api/internal/infrastructure/database"
	"github.com/gabrielcamargo/oficina-api/internal/infrastructure/observability"
)

func main() {
	cfg, err := configs.Load()
	if err != nil {
		slog.Error("failed to load config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger := observability.NewLogger(cfg.ServiceName, cfg.Environment)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	shutdownTracing, err := observability.SetupTracing(ctx, observability.TracingConfig{
		ServiceName: cfg.ServiceName,
		Environment: cfg.Environment,
		Endpoint:    cfg.OtelExporterEndpoint,
		Headers:     cfg.OtelExporterHeaders,
	})
	if err != nil {
		logger.Error("failed to configure tracing", slog.String("error", err.Error()))
		os.Exit(1)
	}

	if err := database.RunMigrations(cfg.DatabaseURL, "migrations"); err != nil {
		logger.Error("failed to run migrations", slog.String("error", err.Error()))
		os.Exit(1)
	}

	db, err := database.NewConnection(cfg.DatabaseURL)
	if err != nil {
		logger.Error("failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	server := api.NewServer(cfg, db)

	go func() {
		logger.Info("server.started",
			slog.String("event", "server.started"),
			slog.String("port", cfg.ServerPort),
			slog.Bool("tracing_enabled", cfg.OtelExporterEndpoint != ""),
		)
		if err := server.Run(); err != nil {
			logger.Error("server stopped unexpectedly", slog.String("error", err.Error()))
			stop()
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("failed to shut down server gracefully", slog.String("error", err.Error()))
	}
	if err := shutdownTracing(shutdownCtx); err != nil {
		logger.Error("failed to flush traces", slog.String("error", err.Error()))
	}
	logger.Info("server.stopped", slog.String("event", "server.stopped"))
}
