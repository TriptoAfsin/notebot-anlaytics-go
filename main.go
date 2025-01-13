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

const DEFAULT_PORT = "3000"

func main() {
	log.Println("⏳ Loading .env file...")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("🔴 Error loading .env file: %s", err)
	}
	log.Println("🟢 .env file loaded")

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
