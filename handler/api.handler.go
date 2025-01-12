package handler

import "github.com/gofiber/fiber/v2"

func ApiHandler(c *fiber.Ctx) error {
	apiStatus := fiber.Map{
		"endPoints": []string{
			"/",
			"/users",
			"/missed",
			"/notes",
			"/games/notebird",
			"/games/notedino",
		},
		"DB_Connection_Status": false,
		"mode":                 "develop",
	}

	return c.JSON(apiStatus)
}

func HealthCheckHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "OK",
		"db":     false,
	})
}
