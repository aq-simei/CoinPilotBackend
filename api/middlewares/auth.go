package middlewares

import (
	"net/http"
	"strings"

	responses "github.com/aq-simei/coin-pilot/internal"
	"github.com/aq-simei/coin-pilot/internal/config/environment"
	"github.com/aq-simei/coin-pilot/internal/config/logger"
	"github.com/aq-simei/coin-pilot/internal/config/security"
	"github.com/gin-gonic/gin"
)

func ApiKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("x-api-key")
		logger.Info("Passing by ApiKeyMiddleware: %v", apiKey)
		if apiKey == "" || apiKey != environment.GetEnv("API_SECRET") {
			logger.Info("API key missing or invalid")
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	}
}

func JwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		logger.Info("Authorization header: %v", authHeader)
		if authHeader == "" {
			logger.Info("Authorization header missing")
			responses.Unauthorized(c, "Authorization header missing")
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			logger.Info("Authorization header does not satisfy Bearer format")
			responses.Unauthorized(c, "Invalid authorization format")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			logger.Info("Token string is empty")
			responses.Unauthorized(c, "Token is missing")
			return
		}

		userID, err := security.ParseJWT(tokenString)
		if err != nil {
			logger.Info("Failed to parse JWT: %v", err)
			responses.Unauthorized(c, err.Error())
			return
		}

		// Store user_id in context for future use
		c.Set("user_id", userID)
		c.Next()
	}
}
