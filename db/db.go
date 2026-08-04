package db

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB connects using the driver selected by DB_DRIVER ("mysql" default, or "postgres").
// This lets us flip the analytics API between the legacy MySQL DB and the notebot v2 Postgres
// (schema `analytics`) with an env change only — and roll back instantly if needed.
func InitDB() *gorm.DB {
	driver := strings.ToLower(os.Getenv("DB_DRIVER"))
	if driver == "" {
		driver = "mysql"
	}

	var dialector gorm.Dialector

	switch driver {
	case "postgres", "postgresql", "pg":
		// Prefer a full DATABASE_URL (postgres://...); fall back to assembling from parts.
		dsn := os.Getenv("DATABASE_URL")
		if dsn == "" {
			dsn = fmt.Sprintf(
				"host=%s user=%s password=%s dbname=%s port=%s sslmode=require search_path=analytics",
				os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_USER_PASS"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"),
			)
		} else {
			// Ensure the analytics schema is on the search_path (and SSL is on for Railway).
			if !strings.Contains(dsn, "search_path=") {
				sep := "?"
				if strings.Contains(dsn, "?") {
					sep = "&"
				}
				dsn += sep + "search_path=analytics"
			}
			if !strings.Contains(dsn, "sslmode=") {
				dsn += "&sslmode=require"
			}
		}
		dialector = postgres.Open(dsn)

	default: // mysql
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			os.Getenv("DB_USER"), os.Getenv("DB_USER_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"),
		)
		dialector = mysql.Open(dsn)
	}

	log.Printf("⏳ Connecting to database (driver=%s)...", driver)

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("🔴 Failed to connect database (driver=%s): %v", driver, err))
	}

	log.Printf("🟢 Connected to database (driver=%s)", driver)
	DB = gormDB
	return gormDB
}
