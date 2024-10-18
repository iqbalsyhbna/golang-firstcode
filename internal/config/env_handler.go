package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ApiURL     string
	AuthToken  string
	DeviceID   string
	PlatformID string
}

func LoadConfig() (*Config, error) {
	// Try to load .env file if exists, but don't return error if file not found
	_ = godotenv.Load()

	// Load configuration from environment variables
	config := &Config{
		ApiURL:     getEnvWithFallback("API_URL", ""),
		AuthToken:  getEnvWithFallback("AUTH_TOKEN", ""),
		DeviceID:   getEnvWithFallback("DEVICE_ID", ""),
		PlatformID: getEnvWithFallback("PLATFORM_ID", ""),
	}

	return config, nil
}

func getEnvWithFallback(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
