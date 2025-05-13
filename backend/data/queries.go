package data

import "github.com/jackc/pgx/v5/pgxpool"

type Queries struct {
	pool *pgxpool.Pool
}

func NewQueries(pool *pgxpool.Pool) *Queries {
	return &Queries{
		pool,
	}
}
