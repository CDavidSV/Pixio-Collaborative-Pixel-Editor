package data

import (
	"context"

	"github.com/CDavidSV/Pixio/types"
	"github.com/CDavidSV/Pixio/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CanvasQueries struct {
	pool *pgxpool.Pool
}

func (s *CanvasQueries) CreateCanvas(title, description, userID string, width, height uint16, data []byte) (types.Canvas, error) {
	query := `INSERT INTO canvases (canvas_id, owner_id, title, description, width, height, data, access_type) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING last_edited_at, created_at`
	canvasID := utils.GenerateID()

	canvas := types.Canvas{
		ID:          canvasID,
		OwnerID:     userID,
		Title:       title,
		Description: description,
		Width:       width,
		Height:      height,
		PixelData:   data,
		AccessType:  types.Restricted,
	}
	if err := s.pool.QueryRow(context.Background(), query, canvasID, userID, title, description, width, height, data, types.Restricted).Scan(&canvas.LastEditedAt, &canvas.CreatedAt); err != nil {
		return canvas, nil
	}

	return canvas, nil
}
