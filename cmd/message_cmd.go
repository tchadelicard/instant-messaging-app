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

// StartMessageService starts the MessageService daemon
func StartMessageService() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize the database
	config.InitDatabase()

	log.Println("Starting UserService daemon...")

	// Setup RabbitMQ connection and channel
	config.SetupRabbitMQ()
	defer config.CleanupRabbitMQ()

	// Declare the direct exchange for registration, login, and user queries
	config.InitDirectRabbitMQExchange("message_direct_exchange")

	// Declare and bind the getUsers queue
	getMessagesQueue := "message_service_get_messages_queue"
	config.InitQueue(getMessagesQueue)
	config.BindQueueToExchange(getMessagesQueue, "user_direct_exchange", "getMessages")

	// Declare the notification exchange
	config.InitDirectRabbitMQExchange("notification_exchange")

	// Create a context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle system signals for shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		sig := <-sigChan
		log.Printf("Received signal: %v. Initiating shutdown...", sig)
		cancel()
	}()

	// Start consuming getUsers requests
	go func() {
		log.Println("Starting consumer for getUsers queue...")
		handlers.ConsumeGetUsersQueue(ctx, getMessagesQueue, "notification_exchange")
	}()

	// Block until context is canceled
	<-ctx.Done()
	log.Println("UserService daemon stopped gracefully.")
}