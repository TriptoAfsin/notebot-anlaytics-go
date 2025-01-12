package handler

import (
	"github.com/TriptoAfsin/notebot-anlaytics-go/config"
	"github.com/TriptoAfsin/notebot-anlaytics-go/db"
	"github.com/gofiber/fiber/v2"
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
	return c.JSON(fiber.Map{
		"status": "OK",
		"db":     false,
	})
}
