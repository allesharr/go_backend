package api

type AuthResponse struct {
	SessionKey string `json:"session_key"`
	UserId     int    `json:"user_id"`
	Login      string `json:"login"`
}
