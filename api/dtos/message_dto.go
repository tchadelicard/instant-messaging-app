package dtos

import "instant-messaging-app/models"

type MessageDTO struct {
	ID        uint   `json:"id"`
	SenderID   uint `json:"sender_id"`
	ReceiverID uint `json:"receiver_id"`
	Content    string `json:"content"`
}

func ToMessageDTO(message models.Message) MessageDTO {
	return MessageDTO{
		ID:        message.ID,
		SenderID:   message.SenderID,
		ReceiverID: message.ReceiverID,
		Content:    message.Content,
	}
}

func ToMessageDTOs(messages []models.Message) []MessageDTO {
	dtos := make([]MessageDTO, len(messages))
	for i, message := range messages {
		dtos[i] = ToMessageDTO(message)
	}
	return dtos
}