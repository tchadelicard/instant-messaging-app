package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"instant-messaging-app/config"
	"instant-messaging-app/user/handlers"

	"github.com/joho/godotenv"
)

func StartUserService() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize the database
	config.InitDatabase()

	log.Println("Starting UserService daemon...")

	// Setup RabbitMQ connection and channel
	config.SetupRabbitMQ()
	defer config.CleanupRabbitMQ()

	// Declare the direct exchange for registration and login
	config.InitDirectRabbitMQExchange("user_direct_exchange")

	// Declare and bind the registration queue
	registrationQueue := "user_service_registration_queue"
	config.InitQueue(registrationQueue)
	config.BindQueueToExchange(registrationQueue, "user_direct_exchange", "registration")

	// Declare and bind the login queue
	loginQueue := "user_service_login_queue"
	config.InitQueue(loginQueue)
	config.BindQueueToExchange(loginQueue, "user_direct_exchange", "login")

	// Declare the notification exchange
	config.InitDirectRabbitMQExchange("notification_exchange")

	// Create a context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())

	// Handle system signals for shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		sig := <-sigChan
		log.Printf("Received signal: %v. Initiating shutdown...", sig)
		cancel()
	}()

	// Start consuming registration requests
	go handlers.ConsumeRegistrationQueue(ctx, registrationQueue, "notification_exchange")

	// Start consuming login requests
	go handlers.ConsumeLoginQueue(ctx, loginQueue, "notification_exchange")

	// Block until context is canceled
	<-ctx.Done()
	log.Println("UserService daemon stopped gracefully.")
}