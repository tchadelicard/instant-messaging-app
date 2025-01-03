package models

import (
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	SenderID   uint `gorm:"not null" json:"sender_id"`
	Sender     User `gorm:"foreignKey:SenderID" json:"sender"` // Relation avec User
	ReceiverID uint `gorm:"not null" json:"receiver_id"`
	Receiver   User `gorm:"foreignKey:ReceiverID" json:"receiver"` // Relation avec User
	Content    string `gorm:"type:text;not null" json:"content"`
}