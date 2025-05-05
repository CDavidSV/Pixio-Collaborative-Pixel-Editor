package types

import (
	"errors"
	"time"
)

// Common
type Map map[string]any
type NullString string

// Errors
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidToken       = errors.New("invalid token")
	ErrSessionNotFound    = errors.New("session not found")
	ErrSessionExpired     = errors.New("session expired")
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
	ID             string     `json:"id"`
	Username       string     `json:"username"`
	Email          string     `json:"email"`
	CreatedAt      time.Time  `json:"created_at"`
	AvatarURL      NullString `json:"avatar_url"`
	HashedPassword string     `json:"-"`
}

type Session struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
	LastAccessed time.Time `json:"last_accessed"`
	RefreshToken string    `json:"refresh_token"`
}

type UserSession struct {
	ID                   string
	UserID               string
	CreatedAt            time.Time
	ExpiresAt            time.Time
	RefreshToken         string
	AccessToken          string
	AccessTokenExpiresAt time.Time
}

type UserLoginDTO struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8,max=50"`
}
