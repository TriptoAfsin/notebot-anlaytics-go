package handler

import "github.com/gofiber/fiber/v2"

func EntertainmentHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "OK",
	})
}
