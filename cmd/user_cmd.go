package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"instant-messaging-app/config"
	"instant-messaging-app/user/handlers"
	"instant-messaging-app/utils"

	"github.com/joho/godotenv"
)

// StartUserService starts the UserService daemon
func StartUserService() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize the database
	config.InitDatabase()

	log.Println("Starting UserService daemon...")

	// Generate a unique queue name for this instance
	instanceQueue := "user_service_" + utils.GenerateUniqueID()

	// Setup RabbitMQ connection and channel
	rabbitConn, rabbitCh := config.SetupRabbitMQ()
	defer rabbitConn.Close()
	defer rabbitCh.Close()

	// Initialize fan-out bindings for the UserService
	config.InitFanOutRabbitMQBindings(rabbitCh, instanceQueue, "user_fanout_exchange")

	// Setup notification exchange for sending registration status
	config.InitDirectRabbitMQExchange(rabbitCh, "notification_exchange")

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
	go handlers.ConsumeRegistrationQueue(ctx, rabbitCh, instanceQueue, "notification_exchange")

	// Block until context is canceled
	<-ctx.Done()
	log.Println("UserService daemon stopped gracefully.")
}