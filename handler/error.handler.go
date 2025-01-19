package handler

import (
	"log"
	"time"

	"github.com/TriptoAfsin/notebot-anlaytics-go/config"
	"github.com/TriptoAfsin/notebot-anlaytics-go/lib/utils"

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

// PostNewError handles creation of new error logs
func PostNewError(db *gorm.DB, appConfig config.AppConfig) fiber.Handler {
	log.Println("🔵 POST: PostNewError handler called")
	return func(c *fiber.Ctx) error {
		if err := utils.ValidateAdminKey(c, appConfig); err != nil {
			return err
		}

		var errorLog ErrorLog
		if err := c.BodyParser(&errorLog); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"status": config.AppMessages.ErrorLog.BadRequest,
			})
		}

		// Validate required fields
		if errorLog.Email == "" || errorLog.Log == "" || errorLog.OS == "" {
			return c.Status(400).JSON(fiber.Map{
				"status": config.AppMessages.ErrorLog.BadRequest,
			})
		}

		if !utils.ValidateEmail(errorLog.Email) {
			return c.Status(400).JSON(fiber.Map{
				"status": config.AppMessages.ErrorLog.InvalidEmail,
			})
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
			return c.Status(500).JSON(fiber.Map{
				"status": config.AppMessages.ErrorLog.OperationUnsuccessful,
			})
		}

		return c.Status(200).JSON(ErrorResponse{
			ErrorInfo: errorLog,
			Status:    config.AppMessages.ErrorLog.LogInsertSuccess,
		})
	}
}

// GetErrorLogs retrieves all error logs
func GetErrorLogs(db *gorm.DB, appConfig config.AppConfig) fiber.Handler {
	log.Println("🟢 GET: GetErrorLogs handler called")
	return func(c *fiber.Ctx) error {
		if err := utils.ValidateAdminKey(c, appConfig); err != nil {
			return err
		}

		var results []map[string]interface{}
		if err := db.Raw("SELECT * FROM app_err_logs").Scan(&results).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status": config.AppMessages.ErrorLog.OperationUnsuccessful,
			})
		}

		return c.Status(200).JSON(ErrorResponse{
			ErrorLogs: results,
			Status:    config.AppMessages.ErrorLog.LogsFetchSuccess,
		})
	}
}

// GetErrorsByEmail retrieves error logs for a specific email
func GetErrorsByEmail(db *gorm.DB, appConfig config.AppConfig) fiber.Handler {
	log.Println("🟢 GET: GetErrorsByEmail handler called")
	return func(c *fiber.Ctx) error {
		if err := utils.ValidateAdminKey(c, appConfig); err != nil {
			return err
		}

		var body struct {
			Email string `json:"email"`
		}

		if err := c.BodyParser(&body); err != nil || body.Email == "" {
			return c.Status(400).JSON(fiber.Map{
				"status": config.AppMessages.ErrorLog.BadRequest,
			})
		}

		var results []map[string]interface{}
		if err := db.Raw("SELECT * FROM app_err_logs WHERE email LIKE ?", body.Email).
			Scan(&results).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status": config.AppMessages.ErrorLog.EmailFetchError,
			})
		}

		return c.Status(200).JSON(ErrorResponse{
			SearchedLogs: results,
			Status:       config.AppMessages.ErrorLog.LogsFetchSuccess,
		})
	}
}
