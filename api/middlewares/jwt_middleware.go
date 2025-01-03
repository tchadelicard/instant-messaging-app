package middlewares

import (
	"instant-messaging-app/utils"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

// Protected is a middleware function that validates the JWT token
func Protected() fiber.Handler {
	return jwtware.New(jwtware.Config{ // Use jwtware here
		SigningKey:   jwtware.SigningKey{Key: []byte(utils.GetJWTSecret())}, // Fetch the secret key
		ErrorHandler: jwtErrorHandler,        // Handle errors for invalid tokens
	})
}

// jwtErrorHandler handles JWT validation errors
func jwtErrorHandler(c *fiber.Ctx, err error) error {
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: " + err.Error(),
		})
	}
	return nil
}
