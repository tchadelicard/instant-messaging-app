package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"instant-messaging-app/api/routes"
	"instant-messaging-app/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func StartWebServer() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found. Using system environment variables.")
	}

	// Connect to the database
	config.InitDatabase()

	// Set up RabbitMQ connection and channel
	config.SetupRabbitMQ()
	defer config.CleanupRabbitMQ()

	// Declare the notification exchange
	config.InitDirectRabbitMQExchange("notification_exchange")

	// Declare the notification broadcast exchange
	config.InitFanoutRabbitMQExchange("notification_broadcast_exchange")

	// Create a context for managing graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize Fiber app
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Set up routes
	routes.SetupRoutes(app, ctx)

	// Channel to listen for OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start the Fiber app in a goroutine
	fiberErrChan := make(chan error, 1)
	go func() {
		port := os.Getenv("APP_PORT")
		if port == "" {
			port = "5000"
		}
		log.Printf("API Gateway running on port %s", port)
		fiberErrChan <- app.Listen(":" + port)
	}()

	select {
	case sig := <-sigChan:
		log.Printf("Received signal: %v. Initiating shutdown...", sig)
	case err := <-fiberErrChan:
		if err != nil {
			log.Printf("Fiber app error: %v", err)
		}
	}

	// Gracefully shutdown Fiber
	if err := app.Shutdown(); err != nil {
		log.Printf("Error during Fiber shutdown: %v", err)
	}

	// Cancel the context to clean up other resources
	cancel()

	log.Println("API Gateway stopped gracefully.")
}