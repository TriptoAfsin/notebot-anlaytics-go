package route

import (
	"github.com/TriptoAfsin/notebot-anlaytics-go/config"
	"github.com/TriptoAfsin/notebot-anlaytics-go/handler"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RouteInit(app *fiber.App, db *gorm.DB) {

	app.Get("/", handler.ApiHandler)

	app.Get("/health", handler.HealthCheckHandler)

	// NoteBird game routes
	app.Post("/games/notebird", handler.PostNoteBirdScore(db))
	app.Get("/games/notebird", handler.GetNoteBirdHof(db))

	// NoteDino game routes
	app.Post("/games/notedino", handler.PostNoteDinoScore(db))
	app.Get("/games/notedino", handler.GetNoteDinoHof(db))

	// Error logging routes
	app.Post("/logs/err", handler.PostNewError(db, config.GetAppConfig()))
	app.Post("/logs/err/email", handler.GetErrorsByEmail(db, config.GetAppConfig()))
	app.Get("/logs/err", handler.GetErrorLogs(db, config.GetAppConfig()))

	// User routes
	app.Post("/user/new", handler.CreateUser(db))
	app.Get("/users/app", handler.GetAllUsers(db))
	app.Post("/users/app/email", handler.GetUsersByEmail(db))
	app.Post("/users/app/batch_dept", handler.GetUsersByDeptAndBatch(db))

}
