package routes

import (
	"encoding/json"
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
		return websocket.New(func(conn *websocket.Conn) {
			// Create a cancellable context
			ctx, cancel := context.WithCancel(context.Background())
			defer func() {
				cancel() // Cancel the context when the WebSocket connection is closed
				conn.Close()
			}()

			// Read the initial message containing the JWT token
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Failed to read WebSocket message:", err)
				conn.WriteMessage(websocket.TextMessage, []byte(`{"type": "error", "message": "Failed to read authentication token"}`))
				return
			}

			var request types.TokenRequest
			if err := json.Unmarshal(message, &request); err != nil {
				log.Println("Failed to unmarshal request")
				conn.WriteMessage(websocket.TextMessage, []byte(`{"type": "error", "message": "Invalid request format"}`))
				return
			}

			userID, err := utils.ValidateJWT(request.Token)
			if err != nil {
				log.Println("Invalid token:", err)
				conn.WriteMessage(websocket.TextMessage, []byte(`{"type": "error", "message": "Invalid authentication token"}`))
				return
			}

			// Acknowledge successful authentication
			conn.WriteMessage(websocket.TextMessage, []byte(`{"type": "auth", "success": true, "message": "Authenticated successfully"}`))

			// Use the user_id as the identifier for the WebSocket
			queueName := utils.GenerateUUID()

			// Declare a queue for the authenticated user
			if _, err := config.RabbitMQCh.QueueDeclare(
				queueName,
				true,
				true,
				false,
				false,
				nil,
			); err != nil {
				log.Printf("Failed to declare queue for user_id %d: %v", userID, err)
				conn.WriteMessage(websocket.TextMessage, []byte(`{"type": "error", "message": "Failed to initialize user queue"}`))
				return
			}

			// Bind the queue to the exchanges
			if err := config.RabbitMQCh.QueueBind(queueName, queueName, "notification_exchange", false, nil); err != nil {
				log.Printf("Failed to bind queue %s to exchange: %v", queueName, err)
				return
			}
			if err := config.RabbitMQCh.QueueBind(queueName, "", "notification_broadcast_exchange", false, nil); err != nil {
				log.Printf("Failed to bind queue %s to exchange: %v", queueName, err)
				return
			}

			// Handle the WebSocket connection
			handlers.HandleWebSocketConnection(conn, queueName, userID, ctx)
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

		if !config.QueueExists(uuid) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Queue not found for UUID",
			})
		}

		return websocket.New(func(conn *websocket.Conn) {
			// Create a cancellable context
			ctx, cancel := context.WithCancel(context.Background())
			defer func() {
				cancel() // Cancel the context when the WebSocket connection is closed
				conn.Close()
			}()

			handlers.HandleWebSocketConnection(conn, uuid, 0, ctx)
		})(c)
	})

	// Protected routes
	api.Get("/messages/:userId", middlewares.Protected(), controllers.GetMessages) // Retrieve messages
	api.Post("/messages/:userId", middlewares.Protected(), controllers.SendMessage) // Send a message
}