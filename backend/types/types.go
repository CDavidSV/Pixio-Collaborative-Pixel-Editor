package types

import "errors"

// Common
type Map map[string]any

// Errors
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type UserSignupDTO struct {
	Username string `validate:"required,min=3,max=20,alphanum"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8,max=50"`
}

type User struct {
	ID             string `json:"id"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	CreatedAt      int    `json:"created_at"`
	AvatarURL      string `json:"avatar_url"`
	HashedPassword string `json:"-"`
}

type Session struct {
	ID           string `json:"id"`
	UserID       string `json:"user_id"`
	CreatedAt    int    `json:"created_at"`
	ExpiresAt    int    `json:"expires_at"`
	LastAccessed int    `json:"last_accessed"`
	RefreshToken string `json:"refresh_token"`
}

type UserSession struct {
	ID                   string
	UserID               string
	CreatedAt            int64
	ExpiresAt            int64
	RefreshToken         string
	AccessToken          string
	AccessTokenExpiresAt int64
}

type UserLoginDTO struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8,max=50"`
}
