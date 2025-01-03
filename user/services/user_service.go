package services

import (
	"errors"
	"log"

	"instant-messaging-app/config"
	"instant-messaging-app/models"

	"golang.org/x/crypto/bcrypt"
)

func ProcessUserRegistration(username string, password string) error {
	// Check if the user already exists
	var existingUser models.User
	if err := config.DB.Where("username = ?", username).First(&existingUser).Error; err == nil {
		return errors.New("username already taken")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	// Create the user
	user := models.User{
		Username: username,
		Password: string(hashedPassword),
	}

	// Save to the database
	return config.DB.Create(&user).Error
}

func NotifyUnauthenticatedClient(uuid, message string) {
	// Logic to notify WebSocket client
	log.Printf("Notifying client %s: %s", uuid, message)
}