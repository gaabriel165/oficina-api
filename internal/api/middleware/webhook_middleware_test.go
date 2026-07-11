package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gabrielcamargo/oficina-api/internal/api/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func newWebhookContext(t *testing.T, secret, header string) (*httptest.ResponseRecorder, *gin.Engine) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	router := gin.New()
	router.Use(middleware.WebhookAuth(secret))
	router.POST("/hook", func(c *gin.Context) { c.Status(http.StatusOK) })

	request := httptest.NewRequest(http.MethodPost, "/hook", nil)
	if header != "" {
		request.Header.Set("X-Webhook-Secret", header)
	}
	router.ServeHTTP(recorder, request)

	return recorder, router
}

func TestWebhookAuth_ShouldAllowWhenSecretMatches(t *testing.T) {
	recorder, _ := newWebhookContext(t, "top-secret", "top-secret")
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestWebhookAuth_ShouldRejectWhenSecretMismatches(t *testing.T) {
	recorder, _ := newWebhookContext(t, "top-secret", "wrong")
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestWebhookAuth_ShouldRejectWhenHeaderMissing(t *testing.T) {
	recorder, _ := newWebhookContext(t, "top-secret", "")
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestWebhookAuth_ShouldRejectWhenSecretNotConfigured(t *testing.T) {
	recorder, _ := newWebhookContext(t, "", "anything")
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}
