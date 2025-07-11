package utils

import (
	"log/slog"
	"os"
	"yaca/pkg/logger"

	"github.com/joho/godotenv"
)

var LoadEnv = loadEnv

func loadEnv() error {
	// Always try to load .env file for backward compatibility
	err := godotenv.Load()
	
	if err != nil {
		// Check if it's because the file doesn't exist
		if _, statErr := os.Stat(".env"); os.IsNotExist(statErr) {
			// Only log if logger is initialized
			if logger.Logger != nil {
				logger.Debug("No .env file found, using system environment variables")
			}
		} else {
			// Some other error occurred
			if logger.Logger != nil {
				logger.Warn("Failed to load .env file",
					slog.String("error", err.Error()))
			}
		}
		// Return the error for backward compatibility
		return err
	}
	
	if logger.Logger != nil {
		logger.Debug(".env file loaded successfully")
		
		// Log that sensitive environment variables are set (without values)
		logger.Debug("Environment check",
			slog.Bool("has_cloudflare_email", os.Getenv("CLOUDFLARE_API_EMAIL") != ""),
			slog.Bool("has_cloudflare_token", os.Getenv("CLOUDFLARE_API_TOKEN") != ""))
	}

	return nil
}
