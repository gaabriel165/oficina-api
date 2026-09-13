package observability_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/gabrielcamargo/oficina-api/internal/infrastructure/observability"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestIDFromContext_ShouldReturnStoredValue(t *testing.T) {
	ctx := observability.WithRequestID(context.Background(), "req-123")

	assert.Equal(t, "req-123", observability.RequestIDFromContext(ctx))
}

func TestRequestIDFromContext_ShouldReturnEmptyWhenAbsent(t *testing.T) {
	assert.Equal(t, "", observability.RequestIDFromContext(context.Background()))
}

func TestEvent_ShouldWriteJSONWithEventNameAndRequestID(t *testing.T) {
	var buffer bytes.Buffer
	previous := slog.Default()
	t.Cleanup(func() { slog.SetDefault(previous) })

	observability.NewLogger("oficina-api", "test")
	logger := slog.New(slog.NewJSONHandler(&buffer, nil))
	slog.SetDefault(logger)

	ctx := observability.WithRequestID(context.Background(), "req-123")
	observability.Event(ctx, "service_order.created", slog.String("order_id", "os-1"))

	var entry map[string]any
	require.NoError(t, json.Unmarshal(buffer.Bytes(), &entry))
	assert.Equal(t, "service_order.created", entry["event"])
	assert.Equal(t, "service_order.created", entry["msg"])
	assert.Equal(t, "os-1", entry["order_id"])
}
