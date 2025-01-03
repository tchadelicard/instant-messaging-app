package dtos

import (
	"instant-messaging-app/models"
)

type UserDTO struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
}

func ToUserDTO(user models.User) UserDTO {
	return UserDTO{
		ID:        user.ID,
		Username:  user.Username,
	}
}

func ToUserDTOs(users []models.User) []UserDTO {
	dtos := make([]UserDTO, len(users))
	for i, user := range users {
		dtos[i] = ToUserDTO(user)
	}
	return dtos
}