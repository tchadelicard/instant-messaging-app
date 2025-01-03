package controllers

import (
	"instant-messaging-app/api/dtos"
	"instant-messaging-app/api/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// GetMessages retrieves messages between the authenticated user and another user
func GetMessages(c *fiber.Ctx) error {
	// Extract the JWT token from the context
	userToken := c.Locals("user").(*jwt.Token)

	// Extract claims from the token
	claims := userToken.Claims.(jwt.MapClaims)

	// Parse user ID from claims (assuming "user_id" is stored in the JWT)
	userID := uint(claims["user_id"].(float64))

	// Get the ID of the user to fetch messages with
	targetUserID, err := strconv.Atoi(c.Params("userId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Fetch messages using the service
	messages, err := services.GetMessagesBetweenUsers(userID, uint(targetUserID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve messages",
		})
	}

	messagesDTOs := dtos.ToMessageDTOs(messages)

	// Return the messages
	return c.JSON(messagesDTOs)
}

// SendMessage envoie un message à un utilisateur
func SendMessage(c *fiber.Ctx) error {
	userId, err := strconv.Atoi(c.Params("userId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID de l'utilisateur invalide",
		})
	}

	userToken := c.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	currentUserId := uint(claims["user_id"].(float64))

	type Request struct {
		Content string `json:"content"`
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bad request",
		})
	}

	message, err := services.CreateMessage(currentUserId, uint(userId), req.Content)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send message",
		})
	}

	messageDTO := dtos.ToMessageDTO(message)

	return c.JSON(messageDTO)
}