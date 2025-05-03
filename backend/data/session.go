package data

import (
	"context"
	"time"

	"github.com/CDavidSV/Pixio/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionQueries struct {
	pool *pgxpool.Pool
}

func (q *SessionQueries) CreateSession(sessionID, userID, refreshToken string, expiresAt time.Time) (types.Session, error) {
	query := `INSERT INTO sessions (id, user_id, refresh_token, expires_at) VALUES ($1, $2, $3, $4) RETURNING id, user_id, refresh_token, expires_at, created_at, last_accessed`

	row := q.pool.QueryRow(context.Background(), query, sessionID, userID, refreshToken, expiresAt)

	var session types.Session
	if err := row.Scan(&session.ID, &session.UserID, &session.RefreshToken, &session.ExpiresAt, &session.CreatedAt, &session.LastAccessed); err != nil {
		return session, err
	}

	return session, nil
}
