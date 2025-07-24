package web

type AuthResponse struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Level int8   `json:"level"`
	Token string `json:"token"`
}
