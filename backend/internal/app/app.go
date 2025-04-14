package app

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	Logger *slog.Logger
	DBPool *pgxpool.Pool
}
