package handlers

import (
	"context"
	"encoding/json"
	"log"

	"instant-messaging-app/config"
	"instant-messaging-app/types"
	"instant-messaging-app/user/services"
	"instant-messaging-app/utils"
)

// ConsumeLoginQueue listens to login requests and processes them
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
				log.Println("Stopping login queue consumption...")
				return
			case msg := <-msgs:
				var request types.AuthenicationRequest
				if err := json.Unmarshal(msg.Body, &request); err != nil {
					log.Printf("Failed to unmarshal login request: %v", err)
					continue
				}

				// Process the login
				success := true
				message := "Login successful"
				token, err := services.ProcessUserLogin(request.Username, request.Password)
				if err != nil {
					success = false
					message = "Login failed: " + err.Error()
					token = ""
				}

				// Publish notification with the message type
				utils.PublishNotification(notificationExchange, request.UUID, "login_response", types.LoginResponse{
					UUID:    request.UUID,
					Success: success,
					Message: message,
					Token:   token,
				})
			}
		}
	}()
}