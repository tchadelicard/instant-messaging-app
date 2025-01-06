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
func HandleWebSocketConnection(conn *websocket.Conn, uuid string, userID uint, ctx context.Context) {
	defer func() {
		conn.Close()

		// Delete the queue when the WebSocket is closed
		if err := config.CleanupQueue(uuid); err != nil {
			log.Printf("Failed to delete queue %s: %v", uuid, err)
		} else {
			log.Printf("Queue %s deleted successfully", uuid)
		}
	}()

	log.Printf("WebSocket connection established for identifier: %s", uuid)

	// Start consuming messages for this WebSocket connection
	go consumeNotifications(ctx, uuid, conn)

	// Keep the WebSocket connection alive
	for {
		_, rawMessage, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket connection closed for identifier: %s", uuid)
			break
		}

		// Handle the incoming message
		if err := handleIncomingWebSocketMessage(conn, rawMessage, uuid, userID); err != nil {
			log.Printf("Failed to handle message for identifier %s: %v", uuid, err)
		}
	}
}

// handleIncomingWebSocketMessage parses and routes the incoming WebSocket message
func handleIncomingWebSocketMessage(conn *websocket.Conn, rawMessage []byte, uuid string, userID uint) error {
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
		if userID == 0 {
			return sendErrorResponse(conn, "Unauthorized request: getUsers requires authentication")
		}
		return handleGetUsers(conn, uuid)
	case "getSelf":
		if userID == 0 {
			return sendErrorResponse(conn, "Unauthorized request: getSelf requires authentication")
		}
		return handleGetSelf(conn, uuid, userID)
	case "getMessages":
		if userID == 0 {
			return sendErrorResponse(conn, "Unauthorized request: getMessages requires authentication")
		}
		return handleGetMessages(conn, uuid, userID, rawMessage)
	default:
		return sendErrorResponse(conn, fmt.Sprintf("Unknown message type: %s", baseMessage.Type))
	}
}

// handleGetUsers retrieves the list of users and sends them to the WebSocket client
func handleGetUsers(conn *websocket.Conn, uuid string) error {
	err := services.PublishGetUsers(uuid)
	if err != nil {
		return sendErrorResponse(conn, fmt.Sprintf("Failed to retrieve users: %v", err))
	}

	return nil
}

func handleGetSelf(conn *websocket.Conn, uuid string, userID uint) error {
	err := services.PublishGetSelf(uuid, userID)
	if err != nil {
		return sendErrorResponse(conn, fmt.Sprintf("Failed to retrieve users: %v", err))
	}

	return nil
}

func handleGetMessages(conn *websocket.Conn, uuid string, userID uint, message []byte) error {
	// Parse the message to extract the recipient ID
	var getMessagesRequest struct {
		Type		string `json:"type"`
		ReceiverID 	uint `json:"receiver_id"`
	}
	json.Unmarshal(message, &getMessagesRequest)

	// Fetch users from the database
	err := services.PublishGetMessages(uuid, userID, getMessagesRequest.ReceiverID)
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
			return err
		}
		return sendMessageToWebSocket(conn, baseMessage)
	case "get_self_response":
		var selfResponse types.GetSelfResponse
		if err := json.Unmarshal(baseMessage.Data, &selfResponse); err != nil {
			return err
		}
		return sendMessageToWebSocket(conn, baseMessage)
	case "get_messages_response":
		var selfResponse types.GetMessagesResponse
		if err := json.Unmarshal(baseMessage.Data, &selfResponse); err != nil {
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