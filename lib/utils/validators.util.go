package utils

import (
	"regexp"

	"github.com/TriptoAfsin/notebot-anlaytics-go/config"
	"github.com/gofiber/fiber/v2"
)

// ValidateAdminKey checks if the provided admin key matches the configured key
func ValidateAdminKey(c *fiber.Ctx, appConfig config.AppConfig) error {
	if adminKey := c.Query("adminKey"); adminKey != appConfig.ADMIN_AUTH_KEY {
		return c.Status(401).JSON(fiber.Map{
			"Error": "🔴 Unauthorized Access !",
		})
	}
	return nil
}

// ValidateEmail checks if the provided email matches a valid email format
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}
