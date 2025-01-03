package routes

import (
	"instant-messaging-app/api/controllers"
	"instant-messaging-app/api/handlers"
	"instant-messaging-app/api/middlewares"
	"instant-messaging-app/config"

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

	// WebSocket route
	app.Get("/ws/:uuid", func(c *fiber.Ctx) error {
		uuid := c.Params("uuid")
		if uuid == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "UUID is required",
			})
		}

		// Check if the queue exists
		if !config.QueueExists(uuid) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Queue not found for UUID",
			})
		}

		// Open the WebSocket connection
		return websocket.New(func(conn *websocket.Conn) {
			handlers.HandleWebSocketConnection(conn, uuid, ctx)
		})(c)
	})

	// Protected routes
	api.Get("/users", middlewares.Protected(), controllers.GetUsers) // Retrieve all users
	api.Get("/users/self", middlewares.Protected(), controllers.GetSelf) // Retrieve the authenticated user
	api.Get("/messages/:userId", middlewares.Protected(), controllers.GetMessages) // Retrieve messages
	api.Post("/messages/:userId", middlewares.Protected(), controllers.SendMessage) // Send a message
}