package handler

import (
	"log"
	"time"

	"github.com/TriptoAfsin/notebot-anlaytics-go/config"
	"github.com/TriptoAfsin/notebot-anlaytics-go/db"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ApiHandler(c *fiber.Ctx) error {
	db := db.DB
	dbStatus := true

	sqlDB, err := db.DB()
	if err != nil || sqlDB.Ping() != nil {
		dbStatus = false
	}

	apiStatus := fiber.Map{
		"endPoints": []string{
			"/",
			"/users",
			"/missed",
			"/notes",
			"/games/notebird",
			"/games/notedino",
		},
		"db_connection": dbStatus,
		"mode":          config.GetAppConfig().ENVIRONMENT,
	}

	return c.JSON(apiStatus)
}

func HealthCheckHandler(c *fiber.Ctx) error {
	db := db.DB
	dbStatus := true

	sqlDB, err := db.DB()
	if err != nil || sqlDB.Ping() != nil {
		dbStatus = false
	}
	return c.JSON(fiber.Map{
		"status": "OK",
		"db":     dbStatus,
	})
}

// GetDailyReport handles fetching daily reports with optional platform filter
func GetDailyReport(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		platform := c.Query("platform")

		var reports []map[string]interface{}
		var query string
		var err error

		if platform != "" {
			query = "SELECT * FROM bot_daily_report WHERE platform = ? ORDER BY date DESC"
			err = db.Raw(query, platform).Scan(&reports).Error
		} else {
			query = "SELECT * FROM bot_daily_report ORDER BY date DESC"
			err = db.Raw(query).Scan(&reports).Error
		}

		if err != nil {
			log.Printf("🔴 Error while fetching daily reports: %v", err)
			return c.Status(500).JSON(fiber.Map{"status": "🔴 Operation was unsuccessful!"})
		}

		return c.Status(200).JSON(fiber.Map{
			"status": "🟢 Operation was successful",
			"data":   reports,
		})
	}
}

// PostDailyReport handles creating or updating daily report entries
func PostDailyReport(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Admin auth check
		if c.Query("adminKey") != config.GetAppConfig().ADMIN_AUTH_KEY {
			return c.Status(401).JSON(fiber.Map{
				"error": "🔴 Unauthorized Access !",
			})
		}

		// Parse request body
		report := struct {
			Platform string `json:"platform"`
		}{}

		if err := c.BodyParser(&report); err != nil {
			return c.Status(400).JSON(fiber.Map{"status": "🔴 Bad Request"})
		}

		// Validate platform
		validPlatforms := []string{"bot", "app"}
		isValidPlatform := false
		for _, p := range validPlatforms {
			if p == report.Platform {
				isValidPlatform = true
				break
			}
		}

		if !isValidPlatform {
			return c.Status(400).JSON(fiber.Map{"status": "🔴 Bad Request - Invalid Platform"})
		}

		currentDate := time.Now().Format("2006-01-02")

		// Check if entry exists
		var count int64
		err := db.Raw("SELECT COUNT(*) FROM bot_daily_report WHERE date = ? AND platform = ?",
			currentDate, report.Platform).Scan(&count).Error

		if err != nil {
			log.Printf("🔴 Error while checking existing log: %v", err)
			return c.Status(500).JSON(fiber.Map{"status": "🔴 Error while checking existing log!"})
		}

		if count > 0 {
			// Update existing entry
			err = db.Exec("UPDATE bot_daily_report SET count = count + 1 WHERE date = ? AND platform = ?",
				currentDate, report.Platform).Error

			if err != nil {
				log.Printf("🔴 Error while updating daily api count: %v", err)
				return c.Status(500).JSON(fiber.Map{"status": "🔴 Error while updating daily api count!"})
			}

			return c.Status(200).JSON(fiber.Map{
				"status": "🟢 Incrementing api call count was successful",
			})
		}

		// Create new entry
		err = db.Exec("INSERT INTO bot_daily_report (date, count, platform) VALUES (?, 1, ?)",
			currentDate, report.Platform).Error

		if err != nil {
			log.Printf("🔴 Error while inserting new log: %v", err)
			return c.Status(500).JSON(fiber.Map{"status": "🔴 Operation was unsuccessful!"})
		}

		return c.Status(200).JSON(fiber.Map{
			"status": "🟢 Creating new log entry was successful",
		})
	}
}
