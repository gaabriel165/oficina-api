package middleware_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gabrielcamargo/oficina-api/internal/api/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestLogger_ShouldLogMethodRouteAndStatus(t *testing.T) {
	var buffer bytes.Buffer
	previous := slog.Default()
	slog.SetDefault(slog.New(slog.NewJSONHandler(&buffer, nil)))
	t.Cleanup(func() { slog.SetDefault(previous) })

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.RequestLogger())
	router.GET("/items/:id", func(c *gin.Context) { c.Status(http.StatusNotFound) })

	router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/items/42", nil))

	var entry map[string]any
	require.NoError(t, json.Unmarshal(buffer.Bytes(), &entry))
	assert.Equal(t, "http.request", entry["event"])
	assert.Equal(t, "GET", entry["method"])
	assert.Equal(t, "/items/:id", entry["route"])
	assert.Equal(t, "/items/42", entry["path"])
	assert.Equal(t, float64(http.StatusNotFound), entry["status"])
	assert.Equal(t, "WARN", entry["level"])
}
