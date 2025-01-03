package controllers

import (
	"instant-messaging-app/api/dtos"
	"instant-messaging-app/api/services"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// GetUsers Controller
func GetUsers(c *fiber.Ctx) error {
	// Retrieve users from the service
	users, err := services.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
	}

	// Convert users to DTOs
	userDTOs := dtos.ToUserDTOs(users)

	// Return the DTOs
	return c.JSON(userDTOs)
}

func GetSelf(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))

	user, err := services.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user",
		})
	}
	userDTO := dtos.ToUserDTO(user)
	return c.JSON(userDTO)
}