package data

import (
	"github.com/CDavidSV/Pixio/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CanvasQueries struct {
	pool *pgxpool.Pool
}

func (s *CanvasQueries) CreateCanvas() types.Canvas {
	return types.Canvas{}
}
