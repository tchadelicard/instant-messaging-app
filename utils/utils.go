package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"instant-messaging-app/config"
	"instant-messaging-app/types"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/google/uuid"

	amqp "github.com/rabbitmq/amqp091-go"
)

func GenerateJWT(user_id uint, username string) (string, error) {
	// Create the Claims
	claims := jwt.MapClaims{
		"user_id": user_id,
		"username": username,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	// Génération du token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(GetJWTSecret())
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return tokenString, nil
}

// getJWTSecret fetches the JWT secret from environment variables
func GetJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return []byte("default_secret")
	}
	return []byte(secret)
}

// GenerateUniqueID creates a unique identifier for this instance
func GenerateUniqueID() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Printf("Failed to get hostname: %v", err)
		hostname = "unknown"
	}
	return strings.ReplaceAll(uuid.New().String()+"_"+hostname, ":", "_")
}

func GenerateUUID() string {
	return uuid.New().String()
}

// validateJWT validates the JWT token and extracts the user ID
func ValidateJWT(tokenString string) (uint, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return GetJWTSecret(), nil
	})
	if err != nil {
		return 0, errors.New("failed to parse token")
	}

	// Extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Retrieve user ID from claims
		userID, ok := claims["user_id"].(float64)
		if !ok {
			return 0, errors.New("invalid claims: user_id not found")
		}
		return uint(userID), nil
	}
	return 0, errors.New("invalid token")
}

func PublishNotification(exchangeName, routingKey, notificationType string, data interface{}) {
	// Create a typed notification
	notification := types.Notification{
		Type: notificationType,
		Data: data,
	}

	// Marshal the notification into JSON
	body, err := json.Marshal(notification)
	if err != nil {
		log.Printf("Failed to marshal notification: %v", err)
		return
	}

	fmt.Println("Publishing message:", string(body))

	// Publish the message to RabbitMQ
	err = config.RabbitMQCh.Publish(
		exchangeName, // Exchange name
		routingKey,   // Routing key
		false,        // Mandatory
		false,        // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Failed to publish notification to exchange %s: %v", exchangeName, err)
	} else {
		log.Printf("Notification published to exchange %s with routing key %s", exchangeName, routingKey)
	}
}