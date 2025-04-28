package data

import "github.com/jackc/pgx/v5/pgxpool"

type Queries struct {
	User UserQueries
}

func NewQueries(pool *pgxpool.Pool) *Queries {
	return &Queries{
		User: UserQueries{pool},
	}
}
