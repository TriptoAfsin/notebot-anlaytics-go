package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	ADMIN_AUTH_KEY string
	ENVIRONMENT    string
}

func GetAppConfig() AppConfig {
	// Load .env file only in non-production environments
	if os.Getenv("ENVIRONMENT") != "production" {
		if err := godotenv.Load(".env"); err != nil {
			log.Printf("⚠️ Warning: Error loading .env file: %s", err)
			// Continue execution, don't fatal
		}
	}

	adminKey := os.Getenv("ADMIN_KEY")
	if adminKey == "" {
		log.Fatal("🔴 ADMIN_KEY environment variable is required")
	}

	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development" // Set default environment
	}

	return AppConfig{
		ADMIN_AUTH_KEY: adminKey,
		ENVIRONMENT:    env,
	}
}
