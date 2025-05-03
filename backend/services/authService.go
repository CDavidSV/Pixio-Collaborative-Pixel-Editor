package services

import (
	"time"

	"github.com/CDavidSV/Pixio/config"
	"github.com/CDavidSV/Pixio/data"
	"github.com/CDavidSV/Pixio/types"
	"github.com/CDavidSV/Pixio/utils"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	queries *data.Queries
}

func (s *AuthService) CreateSession(userID string) (types.UserSession, error) {
	sessionID := utils.GenerateID()

	refreshToken, refreshExpiration, err := s.GenerateRefreshToken(sessionID, userID, config.SessionExpiration)
	if err != nil {
		return types.UserSession{}, err
	}

	accessToken, accessExpiration, err := s.GenerateAccessToken(config.AccessTokenExpiration)
	if err != nil {
		return types.UserSession{}, err
	}

	session, err := s.queries.Session.CreateSession(sessionID, userID, refreshToken, refreshExpiration)
	if err != nil {
		return types.UserSession{}, err
	}

	return types.UserSession{
		ID:                   session.ID,
		UserID:               session.UserID,
		RefreshToken:         refreshToken,
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessExpiration.Unix(),
		ExpiresAt:            refreshExpiration.Unix(),
	}, nil
}

func (s *AuthService) HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(b), err
}

func (s *AuthService) ValidPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *AuthService) GenerateAccessToken(expiresMs int) (string, time.Time, error) {
	expiration := time.Now().Add(time.Duration(expiresMs) * time.Millisecond)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": expiration.Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.AccessTokenSecret))
	if err != nil {
		return "", expiration, err
	}
	return tokenString, expiration, nil
}

func (s *AuthService) GenerateRefreshToken(sessionID, userID string, expiresMs int) (string, time.Time, error) {
	// Implement JWT token generation logic here
	expiration := time.Now().Add(time.Duration(expiresMs) * time.Millisecond)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"session_id": sessionID,
		"user_id":    userID,
		"exp":        expiration.Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.RefreshTokenSecret))
	if err != nil {
		return "", expiration, err
	}
	return tokenString, expiration, nil
}
