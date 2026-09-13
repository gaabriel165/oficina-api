package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gabrielcamargo/oficina-api/internal/api/middleware"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

const testSecret = "test-secret"

func signToken(t *testing.T, claims jwt.MapClaims) string {
	t.Helper()
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(testSecret))
	assert.NoError(t, err)
	return signed
}

func newProtectedRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	protected := router.Group("", middleware.Auth(testSecret))
	protected.GET("/shared", func(c *gin.Context) {
		subject, role := middleware.Principal(c)
		c.JSON(http.StatusOK, gin.H{"subject": subject, "role": role})
	})
	protected.GET("/operators-only", middleware.RequireRole(middleware.RoleOperator), func(c *gin.Context) { c.Status(http.StatusOK) })
	return router
}

func perform(router *gin.Engine, path, token string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodGet, path, nil)
	request.Header.Set("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	return recorder
}

func TestAuth_ShouldDefaultRoleToOperatorWhenClaimMissing(t *testing.T) {
	router := newProtectedRouter()
	token := signToken(t, jwt.MapClaims{"sub": "user-1"})

	recorder := perform(router, "/shared", token)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Body.String(), `"role":"operator"`)
	assert.Contains(t, recorder.Body.String(), `"subject":"user-1"`)
}

func TestRequireRole_ShouldAllowOperator(t *testing.T) {
	router := newProtectedRouter()
	token := signToken(t, jwt.MapClaims{"sub": "user-1", "role": "operator"})

	recorder := perform(router, "/operators-only", token)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestRequireRole_ShouldForbidCustomer(t *testing.T) {
	router := newProtectedRouter()
	token := signToken(t, jwt.MapClaims{"sub": "customer-1", "role": "customer"})

	recorder := perform(router, "/operators-only", token)

	assert.Equal(t, http.StatusForbidden, recorder.Code)
}

func TestRequireRole_ShouldStillAllowCustomerOnSharedRoutes(t *testing.T) {
	router := newProtectedRouter()
	token := signToken(t, jwt.MapClaims{"sub": "customer-1", "role": "customer"})

	recorder := perform(router, "/shared", token)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Body.String(), `"role":"customer"`)
}
