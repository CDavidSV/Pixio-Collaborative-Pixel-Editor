package services

import (
	"errors"
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

	refreshToken, refreshExpiration, err := s.generateRefreshToken(sessionID, userID, config.SessionExpiration)
	if err != nil {
		return types.UserSession{}, err
	}

	session, err := s.queries.Session.CreateSession(sessionID, userID, refreshToken, refreshExpiration)
	if err != nil {
		return types.UserSession{}, err
	}

	accessToken, accessExpiration, err := s.generateAccessToken(config.AccessTokenExpiration, session.UserID)
	if err != nil {
		return types.UserSession{}, err
	}

	return types.UserSession{
		ID:                   session.ID,
		UserID:               session.UserID,
		RefreshToken:         refreshToken,
		AccessToken:          accessToken,
		CreatedAt:            session.CreatedAt,
		AccessTokenExpiresAt: accessExpiration,
		ExpiresAt:            refreshExpiration,
	}, nil
}

func (s *AuthService) HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(b), err
}

func (s *AuthService) ValidPassword(inputPassword, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
	return err == nil
}

func (s *AuthService) CloseSession(refreshToken string) error {
	// Verify refresh token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (any, error) {
		return []byte(config.RefreshTokenSecret), nil
	})

	if err != nil || !token.Valid {
		return types.ErrInvalidToken
	}

	claims := token.Claims.(jwt.MapClaims)
	sessionID, ok := claims["session_id"].(string)
	if !ok {
		return types.ErrInvalidToken
	}

	return s.queries.Session.DeleteSession(sessionID)
}

func (s *AuthService) RevalidateSession(refreshToken string) (types.UserSession, error) {
	// Verify refresh token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (any, error) {
		return []byte(config.RefreshTokenSecret), nil
	})

	if err != nil || !token.Valid {
		return types.UserSession{}, types.ErrInvalidToken
	}

	claims := token.Claims.(jwt.MapClaims)
	sessionID, ok := claims["session_id"].(string)
	if !ok {
		return types.UserSession{}, types.ErrInvalidToken
	}

	session, err := s.queries.Session.GetSession(sessionID)
	if err != nil {
		if errors.Is(err, types.ErrSessionNotFound) {
			// Close all sessions
			s.queries.Session.DeleteAllSessions(session.UserID)
		}

		return types.UserSession{}, err
	}

	// Check if the refresh token matches the one in the database
	if session.RefreshToken != refreshToken {
		// Close all sessions
		s.queries.Session.DeleteAllSessions(session.UserID)
		return types.UserSession{}, types.ErrInvalidToken
	}

	// Check if the session is expired
	if time.Now().After(session.ExpiresAt) {
		return types.UserSession{}, types.ErrSessionExpired
	}

	// Generate new access and refresh tokens
	refreshToken, refreshExpiration, err := s.generateRefreshToken(sessionID, session.UserID, config.SessionExpiration)
	if err != nil {
		return types.UserSession{}, err
	}

	accessToken, accessExpiration, err := s.generateAccessToken(config.AccessTokenExpiration, session.UserID)
	if err != nil {
		return types.UserSession{}, err
	}

	session, err = s.queries.Session.UpdateSession(sessionID, refreshToken, refreshExpiration, time.Now())
	if err != nil {
		return types.UserSession{}, err
	}

	return types.UserSession{
		ID:                   session.ID,
		UserID:               session.UserID,
		RefreshToken:         refreshToken,
		AccessToken:          accessToken,
		CreatedAt:            session.CreatedAt,
		AccessTokenExpiresAt: accessExpiration,
		ExpiresAt:            refreshExpiration,
	}, nil
}

func (s *AuthService) generateAccessToken(expirationTime time.Duration, userID string) (string, time.Time, error) {
	expiration := time.Now().Add(expirationTime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":     expiration.Unix(),
		"user_id": userID,
	})

	tokenString, err := token.SignedString([]byte(config.AccessTokenSecret))
	if err != nil {
		return "", expiration, err
	}
	return tokenString, expiration, nil
}

func (s *AuthService) generateRefreshToken(sessionID, userID string, expirationTime time.Duration) (string, time.Time, error) {
	// Implement JWT token generation logic here
	expiration := time.Now().Add(expirationTime)
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

func (s *AuthService) ValidateAccessToken(accessToken string) (string, bool) {
	// Verify access token
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (any, error) {
		return []byte(config.AccessTokenSecret), nil
	})

	claims := token.Claims.(jwt.MapClaims)
	userID, ok := claims["session_id"].(string)
	if !ok {
		return userID, false
	}

	return userID, err == nil && token.Valid
}
