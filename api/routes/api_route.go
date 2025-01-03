package routes

import (
	"instant-messaging-app/api/controllers"
	"instant-messaging-app/api/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Public routes
	api.Post("/register", controllers.Register)
	api.Post("/login", controllers.Login)

	// Protected routes
	api.Get("/users", middlewares.Protected(), controllers.GetUsers) // Retrieve all users
	api.Get("/users/self", middlewares.Protected(), controllers.GetSelf) // Retrieve the authenticated user
	api.Get("/messages/:userId", middlewares.Protected(), controllers.GetMessages) // Retrieve messages
	api.Post("/messages/:userId", middlewares.Protected(), controllers.SendMessage) // Send a message
}