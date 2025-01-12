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
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("🔴 Error loading .env file: %s", err)
	}
	return AppConfig{
		ADMIN_AUTH_KEY: os.Getenv("ADMIN_KEY"),
		ENVIRONMENT:    os.Getenv("ENVIRONMENT"),
	}
}
