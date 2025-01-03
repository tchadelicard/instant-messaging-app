package utils

import (
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/google/uuid"
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