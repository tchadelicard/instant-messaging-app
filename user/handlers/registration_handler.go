package handlers

import (
	"context"
	"encoding/json"
	"log"

	"instant-messaging-app/types"
	"instant-messaging-app/user/services"

	amqp "github.com/rabbitmq/amqp091-go"
)


func ConsumeRegistrationQueue(ctx context.Context, ch *amqp.Channel, queueName string, notificationExchange string) {
	msgs, _ := ch.Consume(
		queueName, // Queue name
		"",
		true,  // Auto-acknowledge
		false, // Exclusive
		false, // No-local
		false, // No-wait
		nil,
	)

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping registration queue consumption...")
			return
		case msg := <-msgs:
			var request types.RegistrationRequest
			json.Unmarshal(msg.Body, &request)

			// Process the registration
			resultMessage := "Registration successful"
			if err := services.ProcessUserRegistration(request.Username, request.Password); err != nil {
				resultMessage = "Registration failed: " + err.Error()
			}

			// Notify the notification exchange
			notification := types.NotificationMessage{
				UUID:    request.UUID,
				Message: resultMessage,
			}
			body, _ := json.Marshal(notification)

			ch.Publish(
				notificationExchange, // Exchange name
				request.UUID,         // Routing key (UUID as a routing key)
				false,                // Mandatory
				false,                // Immediate
				amqp.Publishing{
					ContentType: "application/json",
					Body:        body,
				},
			)
			log.Printf("Published notification for UUID %s: %s", request.UUID, resultMessage)
		}
	}
}