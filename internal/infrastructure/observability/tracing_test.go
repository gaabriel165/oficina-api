package observability_test

import (
	"context"
	"testing"

	"github.com/gabrielcamargo/oficina-api/internal/infrastructure/observability"
	"github.com/stretchr/testify/assert"
)

func TestSetupTracing_ShouldReturnNoopShutdownWhenDisabled(t *testing.T) {
	shutdown, err := observability.SetupTracing(context.Background(), observability.TracingConfig{})

	assert.NoError(t, err)
	assert.NoError(t, shutdown(context.Background()))
}

func TestSetupTracing_ShouldConfigureExporterWhenEndpointProvided(t *testing.T) {
	cfg := observability.TracingConfig{
		ServiceName: "oficina-api",
		Environment: "test",
		Endpoint:    "http://localhost:4318",
		Headers:     "api-key=secret, other=value",
	}

	shutdown, err := observability.SetupTracing(context.Background(), cfg)

	assert.NoError(t, err)
	assert.NoError(t, shutdown(context.Background()))
}
