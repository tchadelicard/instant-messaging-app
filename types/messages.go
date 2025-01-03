package types

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