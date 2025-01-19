package main

import (
	"fmt"
	"log"
	"os"

	"github.com/TriptoAfsin/notebot-anlaytics-go/db"
	"github.com/TriptoAfsin/notebot-anlaytics-go/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

const DEFAULT_PORT = "10000"

func main() {
	// Try to load .env file, but don't fail if it doesn't exist
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found. Using environment variables...")
	} else {
		log.Println("🟢 .env file loaded")
	}

	// Init DB
	db.InitDB()

	// Init Fiber
	app := fiber.New()

	// Add CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PATCH, DELETE",
	}))

	// Init Route
	routes.RouteInit(app, db.DB)

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = DEFAULT_PORT
	}

	app.Listen(fmt.Sprintf(":%s", port))
}
