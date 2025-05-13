package data

import (
	"context"
	"errors"
	"time"

	"github.com/CDavidSV/Pixio/types"
	"github.com/jackc/pgx/v5"
)

func (q *Queries) CreateSession(sessionID, userID, refreshToken string, expiresAt time.Time) (types.Session, error) {
	query := `INSERT INTO user_sessions (session_id, user_id, refresh_token, expires_at) VALUES ($1, $2, $3, $4) RETURNING session_id, user_id, refresh_token, expires_at, created_at, last_accessed`

	row := q.pool.QueryRow(context.Background(), query, sessionID, userID, refreshToken, expiresAt)

	var session types.Session
	err := row.Scan(&session.ID, &session.UserID, &session.RefreshToken, &session.ExpiresAt, &session.CreatedAt, &session.LastAccessed)
	return session, err
}

func (q *Queries) DeleteSession(sessionID string) error {
	query := `DELETE FROM user_sessions WHERE session_id = $1`

	_, err := q.pool.Exec(context.Background(), query, sessionID)
	return err
}

func (q *Queries) GetSession(sessionID string) (types.Session, error) {
	query := `SELECT session_id, user_id, refresh_token, expires_at, created_at, last_accessed FROM user_sessions WHERE session_id = $1`

	var session types.Session
	if err := q.pool.QueryRow(context.Background(), query, sessionID).Scan(&session.ID, &session.UserID, &session.RefreshToken, &session.ExpiresAt, &session.CreatedAt, &session.LastAccessed); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return session, types.ErrSessionNotFound
		}

		return session, err
	}

	return session, nil
}

func (q *Queries) UpdateSession(sessionID string, refreshToken string, expiresAt time.Time, lastAccessed time.Time) (types.Session, error) {
	query := `UPDATE user_sessions SET refresh_token = $1, expires_at = $2, last_accessed = $3 WHERE session_id = $4 RETURNING session_id, user_id, refresh_token, expires_at, created_at, last_accessed`

	var session types.Session
	err := q.pool.QueryRow(context.Background(), query, refreshToken, expiresAt, lastAccessed, sessionID).Scan(&session.ID, &session.UserID, &session.RefreshToken, &session.ExpiresAt, &session.CreatedAt, &session.LastAccessed)
	return session, err
}

func (q *Queries) DeleteAllSessions(userID string) error {
	query := `DELETE FROM user_sessions WHERE user_id = $1`

	_, err := q.pool.Exec(context.Background(), query, userID)
	return err
}
