package handlers

import (
	"context"
	"encoding/json"
	"instant-messaging-app/api/dtos"
	"instant-messaging-app/config"
	"instant-messaging-app/types"
	"instant-messaging-app/user/services"
	"log"
)

// ConsumeGetUsersQueue listens to getUsers requests and processes them
func ConsumeGetUsersQueue(ctx context.Context, queueName string, notificationExchange string) {
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
				log.Println("Stopping getUsers queue consumption...")
				return
			case msg := <-msgs:
				var request types.GetUsersRequest
				if err := json.Unmarshal(msg.Body, &request); err != nil {
					log.Printf("Failed to unmarshal getUsers request: %v", err)
					continue
				}

				// Fetch users from the database
				users, err := services.GetAllUsers()
				if err != nil {
					log.Printf("Failed to fetch users for user_id %s: %v", request.UserID, err)
					continue
				}

				// Publish notification with the message type
				PublishNotification(notificationExchange, request.UserID, "get_users_response", types.GetUsersResponse{
					Users: dtos.ToUserDTOs(users),
				})

			}
		}
	}()
}