package services

import (
	"errors"
	"instant-messaging-app/config"
	"instant-messaging-app/models"
	"instant-messaging-app/utils"

	"golang.org/x/crypto/bcrypt"
)

// RegisterUser enregistre un utilisateur après validation
func RegisterUser(username, password string) error {
	// Vérifie si l'utilisateur existe déjà
	var existingUser models.User
	if err := config.DB.Where("username = ?", username).First(&existingUser).Error; err == nil {
		return errors.New("username already taken")
	}

	// Hash du mot de passe
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	// Création de l'utilisateur
	user := models.User{
		Username: username,
		Password: string(hashedPassword),
	}

	// Sauvegarde dans la DB
	return config.DB.Create(&user).Error
}

// AuthenticateUser authentifie un utilisateur et retourne un token
func AuthenticateUser(username, password string) (string, error) {
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