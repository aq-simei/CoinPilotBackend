package environment

import (
	"os"

	"github.com/aq-simei/coin-pilot/internal/config/logger"
	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file: %v", err)
	}
}

func GetEnv(key string) string {

	value := os.Getenv(key)
	if value == "" {
		logger.Warn("Environment variable %s is not set", key)
	}
	return value

}
