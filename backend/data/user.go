package data

import (
	"github.com/CDavidSV/Pixio/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserQueries struct {
	pool *pgxpool.Pool
}

func (q *UserQueries) CreateUser(username, email, password string) (types.User, error) {
	return types.User{}, nil
}

func (q *UserQueries) GetUserHashedPassword(email string) (string, error) {
	return "", nil
}
