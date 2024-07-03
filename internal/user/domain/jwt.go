package domain

type Jwt struct {
	ID           string `json:"_id"`
	UserID       string `json:"user_id"`
	RefreshToken string `json:"refresh_token"`
	CreatedAt    string `json:"created_at"`
	Type         string `json:"token_type"`
}
