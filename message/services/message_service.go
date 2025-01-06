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

func CreateMessage(senderID uint, receiverID uint, content string) (models.Message, error) {
	message := models.Message{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
	}

	err := config.DB.Create(&message).Error
	return message, err
}