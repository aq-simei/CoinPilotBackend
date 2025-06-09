package middlewares

import (
	"net/http"
	"strings"

	responses "github.com/aq-simei/coin-pilot/internal"
	"github.com/aq-simei/coin-pilot/internal/config/environment"
	"github.com/aq-simei/coin-pilot/internal/config/logger"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
			responses.InternalServerError(c, "")
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			logger.Info("Authorization header does not satisfy Bearer format")
			responses.InternalServerError(c, "")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			logger.Info("Token string is empty")
			responses.InternalServerError(c, "")
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			jwtSecret := []byte(environment.GetEnv("JWT_SECRET"))
			return jwtSecret, nil
		})

		claims, ok := token.Claims.(jwt.MapClaims)

		// Grab it all, but dont give info on specific error
		if err != nil || !token.Valid || !ok {
			logger.Info("Parsed token claims: %v", claims)
			logger.Info("Token valid: %v", ok)
			logger.Info("Token error: %v", err)
			responses.Unauthorized(c, "Invalid token")
			return
		}

		// Store claims in context for future use
		c.Set("claims", claims)
		c.Next()
	}
}
