package data

import (
	"context"
	"errors"

	"github.com/CDavidSV/Pixio/types"
	"github.com/CDavidSV/Pixio/utils"
	"github.com/jackc/pgx/v5/pgconn"
)

func (q *Queries) CreateUser(username, email, password string) (types.User, error) {
	query := `INSERT INTO users (user_id, username, email, hashed_password) VALUES ($1, $2, $3, $4) RETURNING user_id, username, email, created_at`

	user_id := utils.GenerateID()

	row := q.pool.QueryRow(context.Background(), query, user_id, username, email, password)

	newUser := types.User{}
	if err := row.Scan(&newUser.ID, &newUser.Username, &newUser.Email, &newUser.CreatedAt); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return newUser, types.ErrUserAlreadyExists
		}

		return newUser, err
	}

	return newUser, nil
}

func (q *Queries) GetUserByEmail(email string) (types.User, error) {
	query := `SELECT user_id, username, email, hashed_password, created_at, avatar_url FROM users WHERE email = $1`

	var user types.User
	err := q.pool.QueryRow(context.Background(), query, email).Scan(&user.ID, &user.Username, &user.Email, &user.HashedPassword, &user.CreatedAt, &user.AvatarURL)
	return user, err
}
