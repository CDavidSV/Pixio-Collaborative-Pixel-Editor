package data

import "github.com/jackc/pgx/v5/pgxpool"

type UserQueries struct {
	pool *pgxpool.Pool
}
