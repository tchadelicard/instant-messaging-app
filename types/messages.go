package types

type RegistrationRequest struct {
	UUID     string `json:"uuid"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type NotificationMessage struct {
	UUID    string `json:"uuid"`
	Message string `json:"message"`
}