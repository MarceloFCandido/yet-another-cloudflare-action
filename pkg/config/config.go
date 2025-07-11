package config

import (
	"os"
)

// Config holds the application configuration
type Config struct {
	LogLevel      string
	Debug         bool
	MaskSensitive bool
	Environment   string
}

var AppConfig *Config

// Load initializes the configuration from environment variables
func Load() *Config {
	AppConfig = &Config{
		LogLevel:      getEnvOrDefault("LOG_LEVEL", "INFO"),
		Debug:         getEnvOrDefault("DEBUG", "false") == "true",
		MaskSensitive: getEnvOrDefault("MASK_SENSITIVE", "true") == "true",
		Environment:   getEnvOrDefault("ENVIRONMENT", "development"),
	}
	return AppConfig
}

// getEnvOrDefault returns the environment variable value or a default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
