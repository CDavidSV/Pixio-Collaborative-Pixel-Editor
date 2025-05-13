package types

import (
	"errors"
	"time"
)

// Common
type Map map[string]any
type NullString string
type AccessType int
type AccessRole int
type ObjectType string

const (
	Restricted AccessType = iota
	WithLink
)

const (
	Owner AccessRole = iota
	Editor
	Viewer
)

const (
	CanvasObject     ObjectType = "canvas"
	CollectionObject ObjectType = "collection"
)

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

type Pixel struct {
	R, G, B, A uint8
}

type Canvas struct {
	ID           string     `json:"id"`
	OwnerID      string     `json:"owner_id"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Width        uint16     `json:"width"`
	Height       uint16     `json:"height"`
	PixelData    []byte     `json:"pixel_data"`
	LastEditedAt time.Time  `json:"last_edited_at"`
	AccessType   AccessType `json:"access_type"`
	CreatedAt    time.Time  `json:"created_at"`
	StarCount    uint       `json:"start_count"`
}

type CreateCanvasDTO struct {
	Title       string `json:"title" validate:"required,min=1,max=32"`
	Description string `json:"description" validate:"max=512"`
	Width       uint16 `json:"width" validate:"min=100,max=1024"`
	Height      uint16 `json:"height" validate:"min=100,max=1024"`
}

type DeleteCanvasDTO struct {
	CanvasID string `json:"canvas_id"`
}

type UserAccess struct {
	ObjectID       string
	ObjectType     ObjectType
	UserID         string
	AccessRole     AccessRole
	LastModifiedAt time.Time
	LastModifiedBy string
}
