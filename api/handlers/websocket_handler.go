package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"instant-messaging-app/api/services"
	"instant-messaging-app/config"
	"instant-messaging-app/types"

	"github.com/gofiber/contrib/websocket"
)

// HandleWebSocketConnection manages the WebSocket connection and integrates it with RabbitMQ
func HandleWebSocketConnection(conn *websocket.Conn, identifier string, ctx context.Context, isAuthenticated bool) {
	defer func() {
		conn.Close()

		// Delete the queue when the WebSocket is closed
		if err := config.CleanupQueue(identifier); err != nil {
			log.Printf("Failed to delete queue %s: %v", identifier, err)
		} else {
			log.Printf("Queue %s deleted successfully", identifier)
		}
	}()

	log.Printf("WebSocket connection established for identifier: %s", identifier)

	// Start consuming messages for this WebSocket connection
	go consumeNotifications(ctx, identifier, conn)

	// Keep the WebSocket connection alive
	for {
		_, rawMessage, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket connection closed for identifier: %s", identifier)
			break
		}

		// Handle the incoming message
		if err := handleIncomingWebSocketMessage(conn, rawMessage, identifier, isAuthenticated); err != nil {
			log.Printf("Failed to handle message for identifier %s: %v", identifier, err)
		}
	}
}

// handleIncomingWebSocketMessage parses and routes the incoming WebSocket message
func handleIncomingWebSocketMessage(conn *websocket.Conn, rawMessage []byte, identifier string, isAuthenticated bool) error {
	// Generic message format with a type field
	var baseMessage struct {
		Type string `json:"type"`
	}

	// Parse the message type
	if err := json.Unmarshal(rawMessage, &baseMessage); err != nil {
		return fmt.Errorf("failed to parse message type: %w", err)
	}

	// Route the message based on its type
	switch baseMessage.Type {
	case "getUsers":
		if !isAuthenticated {
			return sendErrorResponse(conn, "Unauthorized request: getUsers requires authentication")
		}
		return handleGetUsers(conn, identifier)
	//case "getMessages":
	//	if !isAuthenticated {
	//		return sendErrorResponse(conn, "Unauthorized request: getMessages requires authentication")
	//	}
	//	return handleGetMessages(conn, rawMessage)
	default:
		return sendErrorResponse(conn, fmt.Sprintf("Unknown message type: %s", baseMessage.Type))
	}
}

// handleGetUsers retrieves the list of users and sends them to the WebSocket client
func handleGetUsers(conn *websocket.Conn, userID string) error {
	// Fetch users from the database
	err := services.PublishGetUsers(userID)
	if err != nil {
		return sendErrorResponse(conn, fmt.Sprintf("Failed to retrieve users: %v", err))
	}

	return nil
}

// sendErrorResponse sends an error response to the WebSocket client
func sendErrorResponse(conn *websocket.Conn, errorMessage string) error {
	response := struct {
		Type  string `json:"type"`
		Error string `json:"error"`
	}{
		Type:  "error",
		Error: errorMessage,
	}

	return sendMessageToWebSocket(conn, response)
}

func consumeNotifications(ctx context.Context, identifier string, conn *websocket.Conn) {
	// Create a new context specifically for this consumer
	consumerCtx, cancel := context.WithCancel(ctx)

	defer func() {
		// Cleanup when the consumer exits
		cancel()
		log.Printf("Consumer cleanup completed for identifier: %s", identifier)

		// Delete the queue after the WebSocket is closed
		if err := config.CleanupQueue(identifier); err != nil {
			log.Printf("Failed to delete queue %s: %v", identifier, err)
		} else {
			log.Printf("Queue %s deleted successfully", identifier)
		}
	}()

	msgs, err := config.RabbitMQCh.Consume(
		identifier, // Queue name
		"",         // Consumer tag
		true,       // Auto-acknowledge
		false,      // Exclusive
		false,      // No-local
		false,      // No-wait
		nil,
	)
	if err != nil {
		log.Printf("Failed to start consumer for queue %s: %v", identifier, err)
		return
	}

	// Listen for messages or context cancellation
	for {
		select {
		case <-consumerCtx.Done():
			log.Printf("Consumer context canceled for identifier: %s", identifier)
			return
		case msg := <-msgs:
			// Process the message
			if err := processMessage(msg.Body, conn); err != nil && len(msg.Body) > 0 {
				log.Printf("Failed to process message for queue %s: %v", identifier, err)
			}
		}
	}
}

// processMessage routes and handles different types of messages
func processMessage(message []byte, conn *websocket.Conn) error {
	// Generic message format with a type field and a data field
	var baseMessage struct {
		Type string          `json:"type"`
		Data json.RawMessage `json:"data"`
	}

	// Parse the base structure of the message
	if err := json.Unmarshal(message, &baseMessage); err != nil {
		return err
	}

	// Route the message based on its type
	switch baseMessage.Type {
	case "registration_response":
		var registrationResponse types.RegistrationResponse
		if err := json.Unmarshal(baseMessage.Data, &registrationResponse); err != nil {
			return err
		}
		return sendMessageToWebSocket(conn, registrationResponse)
	case "login_response":
		var loginResponse types.LoginResponse
		if err := json.Unmarshal(baseMessage.Data, &loginResponse); err != nil {
			return err
		}
		return sendMessageToWebSocket(conn, loginResponse)
	case "get_users_response":
		var usersResponse types.GetUsersResponse
		if err := json.Unmarshal(baseMessage.Data, &usersResponse); err != nil {
			log.Println("I'm here")
			return err
		}
		return sendMessageToWebSocket(conn, baseMessage)
	default:
		log.Printf("Unknown message type: %s", baseMessage.Type)
		return nil
	}
}

// sendMessageToWebSocket sends a structured message to the WebSocket client
func sendMessageToWebSocket(conn *websocket.Conn, message interface{}) error {
	rawMessage, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to serialize message: %w", err)
	}

	return conn.WriteMessage(websocket.TextMessage, rawMessage)
}