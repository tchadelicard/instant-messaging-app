package services

import (
	"instant-messaging-app/config"
	"instant-messaging-app/models"
)

func GetMessagesBetweenUsers(senderID uint, receiverID uint) ([]models.Message, error) {
	var messages []models.Message
	err := config.DB.Preload("Sender").Preload("Receiver").
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			senderID, receiverID, receiverID, senderID).
		Order("created_at asc").
		Find(&messages).Error
	return messages, err
}