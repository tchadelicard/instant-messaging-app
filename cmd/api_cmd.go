package cmd

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

	"instant-messaging-app/api/routes"
	"instant-messaging-app/config"
)

func StartWebServer() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to the database
	config.InitDatabase()

	// Initialize Fiber
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Setup routes
	routes.SetupRoutes(app)

	// Start the server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "5000" // Default port
	}
	log.Printf("Starting web server on port %s", port)
	log.Fatal(app.Listen(":" + port))
}