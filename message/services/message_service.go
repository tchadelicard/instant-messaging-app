package services

import (
	"errors"

	"instant-messaging-app/config"
	"instant-messaging-app/models"
	"instant-messaging-app/utils"

	"golang.org/x/crypto/bcrypt"
)

func ProcessUserRegistration(username, password string) error {
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


// Pr authentifie un utilisateur et retourne un token
func ProcessUserLogin(username, password string) (string, error) {
	var user models.User

	// Vérifie si l'utilisateur existe
	if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return "", errors.New("user not found")
	}

	// Vérifie le mot de passe
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid password")
	}

	// Génère un token JWT
	tokenString, err := utils.GenerateJWT(user.ID, user.Username)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return tokenString, nil
}

// GetAllUsers retrieves all users from the database
func GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := config.DB.Select("id, username").Find(&users).Error
	return users, err
}

func GetUserByID(id uint) (models.User, error) {
	var user models.User
	err := config.DB.First(&user, id).Error
	return user, err
}