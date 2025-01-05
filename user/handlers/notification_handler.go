package handlers

import (
	"encoding/json"
	"fmt"
	"instant-messaging-app/config"
	"instant-messaging-app/types"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// PublishNotification sends a typed notification to the notification exchange
func PublishNotification(exchangeName, routingKey, notificationType string, data interface{}) {
	// Create a typed notification
	notification := types.Notification{
		Type: notificationType,
		Data: data,
	}

	// Marshal the notification into JSON
	body, err := json.Marshal(notification)
	if err != nil {
		log.Printf("Failed to marshal notification: %v", err)
		return
	}

	fmt.Println("Publishing message:", string(body))

	// Publish the message to RabbitMQ
	err = config.RabbitMQCh.Publish(
		exchangeName, // Exchange name
		routingKey,   // Routing key
		false,        // Mandatory
		false,        // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Failed to publish notification to exchange %s: %v", exchangeName, err)
	} else {
		log.Printf("Notification published to exchange %s with routing key %s", exchangeName, routingKey)
	}
}