package handler

import (
	"time"

	"github.com/TriptoAfsin/notebot-anlaytics-go/config"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ErrorLog struct {
	Date  time.Time `json:"date"`
	Log   string    `json:"log"`
	OS    string    `json:"os"`
	Email string    `json:"email"`
}

type ErrorResponse struct {
	ErrorInfo    ErrorLog                 `json:"errorInfo,omitempty"`
	ErrorLogs    []map[string]interface{} `json:"errorLogs,omitempty"`
	SearchedLogs []map[string]interface{} `json:"searched_logs,omitempty"`
	Status       string                   `json:"status"`
}

// validateAdminKey checks if the admin key is valid
func validateAdminKey(c *fiber.Ctx, appConfig config.AppConfig) error {
	if adminKey := c.Query("adminKey"); adminKey != appConfig.ADMIN_AUTH_KEY {
		return c.Status(401).JSON(fiber.Map{
			"Error": "🔴 Unauthorized Access !",
		})
	}
	return nil
}

// PostNewError handles creation of new error logs
func PostNewError(db *gorm.DB, appConfig config.AppConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := validateAdminKey(c, appConfig); err != nil {
			return err
		}

		var errorLog ErrorLog
		if err := c.BodyParser(&errorLog); err != nil {
			return c.Status(400).JSON(fiber.Map{"status": "🔴 Bad Request"})
		}

		// Validate required fields
		if errorLog.Email == "" || errorLog.Log == "" || errorLog.OS == "" {
			return c.Status(400).JSON(fiber.Map{"status": "🔴 Bad Request"})
		}

		if !validateEmail(errorLog.Email) {
			return c.Status(400).JSON(fiber.Map{"status": "🔴 Bad Request, Invalid Email"})
		}

		// If date is not provided, use current time
		if errorLog.Date.IsZero() {
			errorLog.Date = time.Now()
		}

		// Execute raw SQL query
		result := db.Exec(`
			INSERT INTO app_err_logs (date, log, os, email) 
			VALUES (?, ?, ?, ?)`,
			errorLog.Date, errorLog.Log, errorLog.OS, errorLog.Email,
		)

		if result.Error != nil {
			return c.Status(500).JSON(fiber.Map{"status": "🔴 Operation was unsuccessful!"})
		}

		return c.Status(200).JSON(ErrorResponse{
			ErrorInfo: errorLog,
			Status:    "🟢 New Error log insertion was successful",
		})
	}
}

// GetErrorLogs retrieves all error logs
func GetErrorLogs(db *gorm.DB, appConfig config.AppConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := validateAdminKey(c, appConfig); err != nil {
			return err
		}

		var results []map[string]interface{}
		if err := db.Raw("SELECT * FROM app_err_logs").Scan(&results).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "🔴 Operation was unsuccessful!"})
		}

		return c.Status(200).JSON(ErrorResponse{
			ErrorLogs: results,
			Status:    "🟢 Logs Data fetching was successful",
		})
	}
}

// GetErrorsByEmail retrieves error logs for a specific email
func GetErrorsByEmail(db *gorm.DB, appConfig config.AppConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := validateAdminKey(c, appConfig); err != nil {
			return err
		}

		var body struct {
			Email string `json:"email"`
		}

		if err := c.BodyParser(&body); err != nil || body.Email == "" {
			return c.Status(400).JSON(fiber.Map{"status": "🔴 Bad Request"})
		}

		var results []map[string]interface{}
		if err := db.Raw("SELECT * FROM app_err_logs WHERE email LIKE ?", body.Email).
			Scan(&results).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status": "🔴 Error while fetching logs by email",
			})
		}

		return c.Status(200).JSON(ErrorResponse{
			SearchedLogs: results,
			Status:       "🟢 Logs fetching was successful",
		})
	}
}
