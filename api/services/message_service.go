package services

import (
	"encoding/json"
	"fmt"
	"instant-messaging-app/config"
	"instant-messaging-app/models"
	"instant-messaging-app/types"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// GetMessagesBetweenUsers récupère les messages entre deux utilisateurs
func GetMessagesBetweenUsers(senderID uint, receiverID uint) ([]models.Message, error) {
	var messages []models.Message
	err := config.DB.Preload("Sender").Preload("Receiver").
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			senderID, receiverID, receiverID, senderID).
		Order("created_at asc").
		Find(&messages).Error
	return messages, err
}

// CreateMessage crée un nouveau message
func CreateMessage(senderID uint, receiverID uint, content string) (models.Message, error) {
	message := models.Message{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
	}

	err := config.DB.Create(&message).Error
	return message, err
}

func PublishGetMessages(uuid string, userID, receiverID uint) error {
	// Define the registration request payload
	request := types.GetMessagesRequest{
		UUID: uuid,
		UserID: userID,
		ReceiverID: receiverID,
	}

	// Marshal the request to JSON
	body, err := json.Marshal(request)
	if err != nil {
		log.Printf("Failed to marshal GetMessagesRequest: %v", err)
		return fmt.Errorf("failed to marshal GetMessagesRequest")
	}

	// Publish the message to the "user_direct_exchange" with the routing key "registration"
	err = config.RabbitMQCh.Publish(
		"user_direct_exchange", // Exchange name
		"getMessages",         // Routing key
		false,                  // Mandatory
		false,                  // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Failed to publish getMessages request: %v", err)
		return fmt.Errorf("failed to publish getMessages request")
	}

	log.Printf("Published getMessages request for userID %v", userID)
	return nil
}

func PublishSendMessage(uuid string, userID, receiverID uint, content string) error {
	// Define the registration request payload
	request := types.SendMessageRequest{
		UUID: uuid,
		UserID: userID,
		ReceiverID: receiverID,
		Content: content,
	}

	// Marshal the request to JSON
	body, err := json.Marshal(request)
	if err != nil {
		log.Printf("Failed to marshal GetMessagesRequest: %v", err)
		return fmt.Errorf("failed to marshal GetMessagesRequest")
	}

	// Publish the message to the "user_direct_exchange" with the routing key "registration"
	err = config.RabbitMQCh.Publish(
		"user_direct_exchange", // Exchange name
		"sendMessage",         // Routing key
		false,                  // Mandatory
		false,                  // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Failed to publish sendMessage request: %v", err)
		return fmt.Errorf("failed to publish sendMessage request")
	}

	log.Printf("Published sendMessage request for userID %v", userID)
	return nil
}