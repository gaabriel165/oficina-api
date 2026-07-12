package middleware

import (
	"crypto/subtle"
	"net/http"

	"github.com/gin-gonic/gin"
)

const webhookSecretHeader = "X-Webhook-Secret"

func WebhookAuth(webhookSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if webhookSecret == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "webhook authentication not configured"})
			return
		}

		provided := c.GetHeader(webhookSecretHeader)
		if subtle.ConstantTimeCompare([]byte(provided), []byte(webhookSecret)) != 1 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid webhook secret"})
			return
		}

		c.Next()
	}
}
