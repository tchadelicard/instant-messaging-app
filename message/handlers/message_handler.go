package handlers

import (
	"context"
	"encoding/json"
	"instant-messaging-app/config"
	"instant-messaging-app/dtos"
	"instant-messaging-app/message/services"
	"instant-messaging-app/types"
	"instant-messaging-app/utils"
	"log"
)

// ConsumeGetUsersQueue listens to getUsers requests and processes them
func ConsumeGetMessagesQueue(ctx context.Context, queueName string, notificationExchange string) {
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
				log.Println("Stopping getMessages queue consumption...")
				return
			case msg := <-msgs:
				var request types.GetMessagesRequest
				if err := json.Unmarshal(msg.Body, &request); err != nil {
					log.Printf("Failed to unmarshal getMessages request: %v", err)
					continue
				}

				// Fetch users from the database
				log.Printf("Fetching messages between %v and %v", request.UserID, request.ReceiverID)
				messages, err := services.GetMessagesBetweenUsers(request.UserID, request.ReceiverID)
				if err != nil {
					log.Printf("Failed to fetch messages for user id: %s: %v", request.UUID, err)
					continue
				}

				// Publish notification with the message type
				utils.PublishNotification(notificationExchange, request.UUID, "get_messages_response", types.GetMessagesResponse{
					Messages: dtos.ToMessageDTOs(messages),
				})

			}
		}
	}()
}

// ConsumeGetUsersQueue listens to getUsers requests and processes them
func ConsumeSendMessageQueue(ctx context.Context, queueName string, notificationExchange string) {
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
				log.Println("Stopping sendMessage queue consumption...")
				return
			case msg := <-msgs:
				var request types.SendMessageRequest
				if err := json.Unmarshal(msg.Body, &request); err != nil {
					log.Printf("Failed to unmarshal sendMessage request: %v", err)
					continue
				}

				// Fetch users from the database
				log.Printf("Fetching messages between %v and %v", request.UserID, request.ReceiverID)
				message, err := services.CreateMessage(request.UserID, request.ReceiverID, request.Content)
				if err != nil {
					log.Printf("Failed to fetch messages for user id: %s: %v", request.UUID, err)
					continue
				}

				// Publish notification with the message type
				utils.PublishNotification(notificationExchange, "", "send_message_response", types.SendMessageResponse{
					Message: dtos.ToMessageDTO(message),
				})
			}
		}
	}()
}