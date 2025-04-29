package services

import (
	"github.com/CDavidSV/Pixio/data"
	"github.com/CDavidSV/Pixio/types"
)

type AuthService struct {
	queries *data.Queries
}

func (s *AuthService) Authenticate(email, password string) error {
	return nil
}

func (s *AuthService) CreateSession(userID string) (types.UserSession, error) {
	return types.UserSession{}, nil
}
