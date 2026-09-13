package middleware

import (
	"github.com/gabrielcamargo/oficina-api/internal/infrastructure/observability"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const RequestIDHeader = "X-Request-ID"

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.NewString()
		}

		c.Header(RequestIDHeader, requestID)
		c.Request = c.Request.WithContext(observability.WithRequestID(c.Request.Context(), requestID))
		c.Next()
	}
}
