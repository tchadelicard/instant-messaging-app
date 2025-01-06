package routes

import (
	"encoding/json"
	"fmt"
	"instant-messaging-app/api/controllers"
	"instant-messaging-app/api/handlers"
	"instant-messaging-app/api/middlewares"
	"instant-messaging-app/config"
	"instant-messaging-app/types"
	"instant-messaging-app/utils"
	"log"

	"context"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures all application routes
func SetupRoutes(app *fiber.App, ctx context.Context) {
	api := app.Group("/api")

	// Public routes
	api.Post("/register", controllers.Register)
	api.Post("/login", controllers.Login)

	// Authenticated WebSocket route (for chat and other interactions)
	app.Get("/ws/auth", func(c *fiber.Ctx) error {
		// Open the WebSocket connection
		return websocket.New(func(conn *websocket.Conn) {
			defer conn.Close()

			// Read the initial message containing the JWT token
			_, message, err := conn.ReadMessage()
			log.Println(string(message))
			if err != nil {
				log.Println("Failed to read WebSocket message:", err)
				conn.WriteMessage(websocket.TextMessage, []byte(`{"type": "error", "message": "Failed to read authentication token"}`))
				return
			}

			var request types.TokenRequest
			err = json.Unmarshal(message, &request)
			if (err != nil) {
				log.Printf("Failed to unmarshal request")
				return
			}

			userID, err := utils.ValidateJWT(request.Token)
			if err != nil {
				log.Printf("Invalid token: %v", err)
				conn.WriteMessage(websocket.TextMessage, []byte(`{"type": "error", "message": "Invalid authentication token"}`))
				return
			}

			// Acknowledge successful authentication
			conn.WriteMessage(websocket.TextMessage, []byte(`{"type": "auth", "success": true, "message": "Authenticated successfully"}`))

			// Use the user_id as the identifier for the WebSocket
			queueName := utils.GenerateUUID()

			// Declare a queue for the authenticated user
			_, err = config.RabbitMQCh.QueueDeclare(
				queueName, // Queue name
				true,      // Durable
				true,      // Auto-delete (deleted when last consumer disconnects)
				false,     // Exclusive
				false,     // No-wait
				nil,       // Arguments
			)
			if err != nil {
				log.Printf("Failed to declare queue for user_id %d: %v", userID, err)
				conn.WriteMessage(websocket.TextMessage, []byte(`{"type": "error", "message": "Failed to initialize user queue"}`))
				return
			}
			// Bind the queue to the notification exchange with the UUID as the routing key
			err = config.RabbitMQCh.QueueBind(
				queueName,             // Queue name
				queueName,             // Routing key
				"notification_exchange", // Exchange name
				false,
				nil,
			)
			if err != nil {
				log.Printf("Failed to bind queue %s to exchange: %v", queueName, err)
				return
			}

			fmt.Println("Crash here.")

			// Handle the WebSocket connection
			handlers.HandleWebSocketConnection(conn, queueName, ctx, true)
		})(c)
	})
	// Non-authenticated WebSocket route (for registration and login)
	app.Get("/ws/:uuid", func(c *fiber.Ctx) error {
		uuid := c.Params("uuid")
		if uuid == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "UUID is required",
			})
		}

		log.Println("I'm here.")

		// Check if the queue exists
		if !config.QueueExists(uuid) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Queue not found for UUID",
			})
		}

		// Open the WebSocket connection
		return websocket.New(func(conn *websocket.Conn) {
			handlers.HandleWebSocketConnection(conn, uuid, ctx, false)
		})(c)
	})

	// Protected routes
	api.Get("/users", middlewares.Protected(), controllers.GetUsers)             // Retrieve all users
	api.Get("/users/self", middlewares.Protected(), controllers.GetSelf)         // Retrieve the authenticated user
	api.Get("/messages/:userId", middlewares.Protected(), controllers.GetMessages) // Retrieve messages
	api.Post("/messages/:userId", middlewares.Protected(), controllers.SendMessage) // Send a message
}