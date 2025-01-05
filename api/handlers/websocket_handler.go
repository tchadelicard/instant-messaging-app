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
	defer conn.Close()

	// Log the connection type
	if isAuthenticated {
		log.Printf("Authenticated WebSocket connection established for user_id: %s", identifier)
	} else {
		log.Printf("Unauthenticated WebSocket connection established for UUID: %s", identifier)
	}

	// Start consuming messages for the identifier (user_id or UUID)
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
			// Process the message based on its type
			log.Println(string(msg.Body))
			if err := processMessage(msg.Body, conn); err != nil {
				log.Printf("Failed to process message for queue %s: %v", queueName, err)
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