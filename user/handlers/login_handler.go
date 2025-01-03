package handlers

import (
	"context"
	"encoding/json"
	"log"

	"instant-messaging-app/config"
	"instant-messaging-app/types"
	"instant-messaging-app/user/services"

	amqp "github.com/rabbitmq/amqp091-go"
)

// ConsumeRegistrationQueue listens to registration requests and processes them
func ConsumeLoginQueue(ctx context.Context, queueName string, notificationExchange string) {
	msgs, err := config.RabbitMQCh.Consume(
		queueName, // Queue name
		"",
		true,  // Auto-acknowledge
		false, // Exclusive
		false, // No-local
		false, // No-wait
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to start consuming from queue %s: %v", queueName, err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("Stopping registration queue consumption...")
				return
			case msg := <-msgs:
				var request types.AuthenicationRequest
				if err := json.Unmarshal(msg.Body, &request); err != nil {
					log.Printf("Failed to unmarshal registration request: %v", err)
					continue
				}

				// Process the registration
				success := true
				message := "Login successful"
				token, err := services.ProcessUserLogin(request.Username, request.Password)
				if err != nil {
					success = false
					message = "Login failed: " + err.Error()
				}

				// Publish notification with UUID as the routing key
				publishLoginNotification(notificationExchange, request.UUID, success, message, token)
			}
		}
	}()
}

// publishLoginNotification sends the notification to the notification exchange
func publishLoginNotification(exchangeName, uuid string, success bool, message, token string) {
	notification := types.LoginResponse{
		UUID:    uuid,
		Success: success,
		Message: message,
		Token:   token,
	}

	body, err := json.Marshal(notification)
	if err != nil {
		log.Printf("Failed to marshal notification for UUID %s: %v", uuid, err)
		return
	}

	err = config.RabbitMQCh.Publish(
		exchangeName, // Exchange name
		uuid,         // Routing key
		false,        // Mandatory
		false,        // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Failed to publish notification for UUID %s: %v", uuid, err)
	} else {
		log.Printf("Notification published for UUID %s: %s", uuid, message)
	}
}