package route

import (
	"github.com/TriptoAfsin/notebot-anlaytics-go/handler"
	"github.com/gofiber/fiber/v2"
)

func RouteInit(app *fiber.App) {

	app.Get("/", handler.ApiHandler)

	app.Get("/health", handler.HealthCheckHandler)
}
