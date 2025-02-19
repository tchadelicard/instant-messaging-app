package types

import (
	"instant-messaging-app/dtos"
)

type AuthenicationRequest struct {
	UUID     string `json:"uuid"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegistrationResponse struct {
	UUID    string `json:"uuid"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type LoginResponse struct {
	UUID	string `json:"uuid"`
	Success	bool   `json:"success"`
	Message string `json:"message"`
	Token	string `json:"token"`
}

// Notification represents the generic notification structure
type Notification struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type GetUsersRequest struct {
	UUID 	string	`json:"uuid"`
}

type GetUsersResponse struct {
	Users	[]dtos.UserDTO	`json:"users"`
}

type TokenRequest struct {
	Type	string	`json:"type"`
	Token	string	`json:"token"`
}

type GetSelfRequest struct {
	UUID 	string	`json:"uuid"`
	UserID 	uint	`json:"user_id"`
}

type GetSelfResponse struct {
	User	dtos.UserDTO	`json:"user"`
}

type GetMessagesRequest struct {
	UUID 		string	`json:"uuid"`
	UserID		uint	`json:"user_id"`
	ReceiverID	uint	`json:"receiver_id"`
}

type GetMessagesResponse struct {
	Messages	[]dtos.MessageDTO	`json:"messages"`
}

type SendMessageRequest struct {
	UUID 		string	`json:"uuid"`
	UserID		uint	`json:"user_id"`
	ReceiverID	uint	`json:"receiver_id"`
	Content		string	`json:"content"`
}

type SendMessageResponse struct {
	Message	dtos.MessageDTO	`json:"message"`
}