package security

import (
	"errors"
	"time"

	"github.com/aq-simei/coin-pilot/internal/config/environment"
	"github.com/aq-simei/coin-pilot/internal/config/logger"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := environment.GetEnv("JWT_SECRET")
	return token.SignedString([]byte(secret))
}

func ParseJWT(tokenString string) (string, error) {
	secret := environment.GetEnv("JWT_SECRET")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is as expected
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return "", err
	}

	// Extract claims and validate
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(string)
		if !ok || userID == "" {
			return "", errors.New("user_id claim not found")
		}
		if exp, ok := claims["exp"].(float64); ok {
			if time.Unix(int64(exp), 0).Before(time.Now()) {
				return "", errors.New("token has expired")
			}
		}
		logger.Info("JWT parsed successfully, user_id: %s", userID)
		logger.Info("JWT expiration time: %v", time.Unix(int64(claims["exp"].(float64)), 0))
		return userID, nil
	}

	return "", errors.New("invalid token")
}
