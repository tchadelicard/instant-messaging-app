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


func PublishGetUsers(uuid string) error {
	// Define the registration request payload
	request := types.GetUsersRequest{
		UUID: uuid,
	}

	// Marshal the request to JSON
	body, err := json.Marshal(request)
	if err != nil {
		log.Printf("Failed to marshal registration request: %v", err)
		return fmt.Errorf("failed to marshal registration request")
	}

	// Publish the message to the "user_direct_exchange" with the routing key "registration"
	err = config.RabbitMQCh.Publish(
		"user_direct_exchange", // Exchange name
		"getUsers",         // Routing key
		false,                  // Mandatory
		false,                  // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Failed to publish getUsers request: %v", err)
		return fmt.Errorf("failed to publish getUsers request")
	}

	log.Printf("Published getUsers request for userID %s", uuid)
	return nil
}

func PublishGetSelf(uuid string, userID uint) error {
	// Define the registration request payload
	request := types.GetSelfRequest{
		UUID: uuid,
		UserID: userID,
	}

	// Marshal the request to JSON
	body, err := json.Marshal(request)
	if err != nil {
		log.Printf("Failed to marshal registration request: %v", err)
		return fmt.Errorf("failed to marshal registration request")
	}

	// Publish the message to the "user_direct_exchange" with the routing key "registration"
	err = config.RabbitMQCh.Publish(
		"user_direct_exchange", // Exchange name
		"getSelf",         // Routing key
		false,                  // Mandatory
		false,                  // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Failed to publish getSelf request: %v", err)
		return fmt.Errorf("failed to publish getSelf request")
	}

	log.Printf("Published getSelf request for userID %s", uuid)
	return nil
}

// GetAllUsers retrieves all users from the database
func GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := config.DB.Select("id, username").Find(&users).Error
	return users, err
}

func GetUserByID(id uint) (models.User, error) {
	var user models.User
	err := config.DB.First(&user, id).Error
	return user, err
}