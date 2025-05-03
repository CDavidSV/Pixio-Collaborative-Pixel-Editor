package data

import "github.com/jackc/pgx/v5/pgxpool"

type Queries struct {
	User    UserQueries
	Session SessionQueries
}

func NewQueries(pool *pgxpool.Pool) *Queries {
	return &Queries{
		User:    UserQueries{pool},
		Session: SessionQueries{pool},
	}
}
