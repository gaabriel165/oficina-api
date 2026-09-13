package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	RoleOperator = "operator"
	RoleCustomer = "customer"

	contextKeySubject = "user_id"
	contextKeyRole    = "role"
)

func Auth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid authorization header"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}

		c.Set(contextKeySubject, claims["sub"])
		c.Set(contextKeyRole, roleFromClaims(claims))
		c.Next()
	}
}

func RequireRole(allowed ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, role := Principal(c)
		for _, candidate := range allowed {
			if role == candidate {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "this operation is not allowed for your role"})
	}
}

func Principal(c *gin.Context) (subject, role string) {
	subject, _ = c.Value(contextKeySubject).(string)
	role, _ = c.Value(contextKeyRole).(string)
	return subject, role
}

func roleFromClaims(claims jwt.MapClaims) string {
	role, ok := claims["role"].(string)
	if !ok || role == "" {
		return RoleOperator
	}
	return role
}
