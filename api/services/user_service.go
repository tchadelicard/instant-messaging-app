package services

import (
	"instant-messaging-app/config"
	"instant-messaging-app/models"
)

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