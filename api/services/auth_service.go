package services

import (
	"encoding/json"
	"fmt"
	"instant-messaging-app/config"
	"instant-messaging-app/types"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// PublishRegistrationRequest publishes a registration request to RabbitMQ
func PublishRegistrationRequest(uuid, username, password string) error {
	// Define the registration request payload
	request := types.AuthenicationRequest{
		UUID:     uuid,
		Username: username,
		Password: password,
	}

	// Marshal the request to JSON
	body, err := json.Marshal(request)
	if err != nil {
		log.Printf("Failed to marshal registration request: %v", err)
		return fmt.Errorf("failed to marshal registration request")
	}
	// Create and bind a queue for the UUID
	_, err = config.RabbitMQCh.QueueDeclare(
		uuid, // Queue name
		true,      // Durable
		true,      // Auto-delete (deleted when last consumer disconnects)
		false,     // Exclusive
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		log.Printf("Failed to declare queue for UUID %s: %v", uuid, err)
		return fmt.Errorf("failed to declare queue for UUID %s: %w", uuid, err)
	}

	// Bind the queue to the notification exchange with the UUID as the routing key
	err = config.RabbitMQCh.QueueBind(
		uuid,             // Queue name
		uuid,             // Routing key
		"notification_exchange", // Exchange name
		false,
		nil,
	)
	if err != nil {
		log.Printf("Failed to bind queue %s to exchange: %v", uuid, err)
		return fmt.Errorf("failed to bind queue %s: %w", uuid, err)
	}

	log.Printf("Queue created and bound for UUID: %s", uuid)

	// Publish the message to the "user_direct_exchange" with the routing key "registration"
	err = config.RabbitMQCh.Publish(
		"user_direct_exchange", // Exchange name
		"registration",         // Routing key
		false,                  // Mandatory
		false,                  // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Failed to publish registration request: %v", err)
		return fmt.Errorf("failed to publish registration request")
	}

	log.Printf("Published registration request for UUID %s", uuid)
	return nil
}

// PublishRegistrationRequest publishes a registration request to RabbitMQ
func PublishLoginRequest(uuid, username, password string) error {
	// Define the registration request payload
	request := types.AuthenicationRequest{
		UUID:     uuid,
		Username: username,
		Password: password,
	}

	// Marshal the request to JSON
	body, err := json.Marshal(request)
	if err != nil {
		log.Printf("Failed to marshal registration request: %v", err)
		return fmt.Errorf("failed to marshal registration request")
	}
	// Create and bind a queue for the UUID
	_, err = config.RabbitMQCh.QueueDeclare(
		uuid, // Queue name
		true,      // Durable
		true,      // Auto-delete (deleted when last consumer disconnects)
		false,     // Exclusive
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		log.Printf("Failed to declare queue for UUID %s: %v", uuid, err)
		return fmt.Errorf("failed to declare queue for UUID %s: %w", uuid, err)
	}

	// Bind the queue to the notification exchange with the UUID as the routing key
	err = config.RabbitMQCh.QueueBind(
		uuid,             // Queue name
		uuid,             // Routing key
		"notification_exchange", // Exchange name
		false,
		nil,
	)
	if err != nil {
		log.Printf("Failed to bind queue %s to exchange: %v", uuid, err)
		return fmt.Errorf("failed to bind queue %s: %w", uuid, err)
	}

	log.Printf("Queue created and bound for UUID: %s", uuid)

	// Publish the message to the "user_direct_exchange" with the routing key "registration"
	err = config.RabbitMQCh.Publish(
		"user_direct_exchange", // Exchange name
		"login",         // Routing key
		false,                  // Mandatory
		false,                  // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Failed to publish login request: %v", err)
		return fmt.Errorf("failed to login registration request")
	}

	log.Printf("Published login request for UUID %s", uuid)
	return nil
}