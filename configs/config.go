package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort           string
	Environment          string
	ServiceName          string
	DatabaseURL          string
	JWTSecret            string
	WebhookSecret        string
	ResendAPIKey         string
	EmailFrom            string
	OtelExporterEndpoint string
	OtelExporterHeaders  string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	return &Config{
		ServerPort:           getEnv("SERVER_PORT", "8080"),
		Environment:          getEnv("APP_ENV", "local"),
		ServiceName:          getEnv("OTEL_SERVICE_NAME", "oficina-api"),
		DatabaseURL:          getEnv("DATABASE_URL", ""),
		JWTSecret:            getEnv("JWT_SECRET", ""),
		WebhookSecret:        getEnv("WEBHOOK_SECRET", ""),
		ResendAPIKey:         getEnv("RESEND_API_KEY", ""),
		EmailFrom:            getEnv("EMAIL_FROM", "onboarding@resend.dev"),
		OtelExporterEndpoint: getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", ""),
		OtelExporterHeaders:  getEnv("OTEL_EXPORTER_OTLP_HEADERS", ""),
	}, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
