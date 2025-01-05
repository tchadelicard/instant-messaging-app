package handlers

import (
	"context"
	"encoding/json"
	"log"

	"instant-messaging-app/config"
	"instant-messaging-app/types"
	"instant-messaging-app/user/services"
)

// ConsumeRegistrationQueue listens to registration requests and processes them
func ConsumeRegistrationQueue(ctx context.Context, queueName string, notificationExchange string) {
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
				message := "Registration successful"
				if err := services.ProcessUserRegistration(request.Username, request.Password); err != nil {
					success = false
					message = "Registration failed: " + err.Error()
				}

				// Publish notification with the message type
				PublishNotification(notificationExchange, request.UUID, "registration_response", types.RegistrationResponse{
					UUID:    request.UUID,
					Success: success,
					Message: message,
				})
			}
		}
	}()
}