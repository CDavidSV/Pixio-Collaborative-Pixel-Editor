package data

import "github.com/jackc/pgx/v5/pgxpool"

type CanvasQueries struct {
	pool *pgxpool.Pool
}
