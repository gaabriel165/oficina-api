package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gabrielcamargo/oficina-api/internal/api/middleware"
	"github.com/gabrielcamargo/oficina-api/internal/infrastructure/observability"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRequestID_ShouldPropagateIncomingHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.RequestID())
	var seen string
	router.GET("/ping", func(c *gin.Context) {
		seen = observability.RequestIDFromContext(c.Request.Context())
		c.Status(http.StatusOK)
	})

	request := httptest.NewRequest(http.MethodGet, "/ping", nil)
	request.Header.Set(middleware.RequestIDHeader, "abc-123")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	assert.Equal(t, "abc-123", seen)
	assert.Equal(t, "abc-123", recorder.Header().Get(middleware.RequestIDHeader))
}

func TestRequestID_ShouldGenerateWhenHeaderMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.RequestID())
	router.GET("/ping", func(c *gin.Context) { c.Status(http.StatusOK) })

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/ping", nil))

	assert.NotEmpty(t, recorder.Header().Get(middleware.RequestIDHeader))
}
