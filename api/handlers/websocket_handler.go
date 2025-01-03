package handlers

import (
	"context"
	"fmt"
	"log"

	"instant-messaging-app/config"

	"github.com/gofiber/contrib/websocket"
)

// HandleWebSocketConnection manages the WebSocket connection and integrates it with RabbitMQ
func HandleWebSocketConnection(conn *websocket.Conn, uuid string, ctx context.Context) {
	defer conn.Close()

	log.Printf("WebSocket connection established for UUID: %s", uuid)

	// Start consuming notifications for the given UUID
	go consumeNotifications(ctx, uuid, conn)

	// Keep the WebSocket connection open
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket connection closed for UUID: %s", uuid)
			break
		}
		fmt.Printf("WebSocket message received for UUID: %s, %s\n", uuid, string(message))
	}
}

// consumeNotifications listens to RabbitMQ messages and sends them to the WebSocket client
func consumeNotifications(ctx context.Context, queueName string, conn *websocket.Conn) {
	msgs, err := config.RabbitMQCh.Consume(
		queueName, // Queue name
		"",        // Consumer tag
		true,      // Auto-acknowledge
		false,     // Exclusive
		false,     // No-local
		false,     // No-wait
		nil,
	)
	if err != nil {
		log.Printf("Failed to start consumer for queue %s: %v", queueName, err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			log.Printf("Stopping consumer for queue %s", queueName)
			return
		case msg := <-msgs:
			// Forward the message to the WebSocket client
			err := conn.WriteMessage(websocket.TextMessage, msg.Body)
			if err != nil {
				log.Printf("Failed to send message to WebSocket for queue %s: %v", queueName, err)
				return
			}
		}
	}
}