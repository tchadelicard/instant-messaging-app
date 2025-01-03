package controllers

import (
	"instant-messaging-app/api/services"
	"instant-messaging-app/utils"

	"github.com/gofiber/fiber/v2"
)

// Register Controller
func Register(c *fiber.Ctx) error {
	// Define the request body structure
	type Request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// Parse the request body
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Basic validation
	if req.Username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username is required",
		})
	}
	if len(req.Password) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password must be at least 6 characters long",
		})
	}

	// Generate a UUID for WebSocket tracking
	uuid := utils.GenerateUUID()

	// Publish the registration request to RabbitMQ
	err := services.PublishRegistrationRequest(uuid, req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process registration",
		})
	}

	// Respond with the UUID for WebSocket tracking
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"uuid": uuid,
		"message": "Registration request received. Use the UUID to track status via WebSocket.",
	})
}

// Login Controller
func Login(c *fiber.Ctx) error {
	// Define the request body structure
	type Request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// Parse the request body
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Basic validation
	if req.Username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username is required",
		})
	}
	if req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password is required",
		})
	}

	// Authenticate the user and get a token
	token, err := services.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		if err.Error() == "invalid password" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid password",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	// Success response with token
	return c.JSON(fiber.Map{
		"token": token,
	})
}