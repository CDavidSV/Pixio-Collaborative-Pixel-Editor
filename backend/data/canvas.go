package data

import (
	"context"
	"fmt"
	"strings"

	"github.com/CDavidSV/Pixio/types"
	"github.com/CDavidSV/Pixio/utils"
	"github.com/jackc/pgx/v5"
)

func (q *Queries) CreateCanvas(title, description, userID string, width, height uint16, data []byte) (types.Canvas, error) {
	canvasID := utils.GenerateID()
	var canvas types.Canvas

	ctx := context.Background()
	tx, err := q.pool.Begin(ctx)
	if err != nil {
		return canvas, err
	}
	defer tx.Rollback(ctx)

	batch := &pgx.Batch{}

	// Insert canvas
	batch.Queue(`
		INSERT INTO canvases (canvas_id, owner_id, title, description, width, height, data, access_type)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING last_edited_at, created_at;
	`, canvasID, userID, title, description, width, height, data, types.Restricted)

	// Insert access rule
	batch.Queue(`
		INSERT INTO access_rules (object_id, object_type, user_id, permissions)
		VALUES ($1, 'canvas', $2, $3);
	`, canvasID, userID, types.Owner)

	br := tx.SendBatch(ctx, batch)
	defer br.Close()

	canvas = types.Canvas{
		ID:             canvasID,
		OwnerID:        userID,
		Title:          title,
		Description:    description,
		Width:          width,
		Height:         height,
		PixelData:      data,
		LinkAccessType: types.Restricted,
		LinkAccessRole: types.Viewer,
	}

	if err := br.QueryRow().Scan(&canvas.LastEditedAt, &canvas.CreatedAt); err != nil {
		return canvas, fmt.Errorf("failed to insert canvas: %w", err)
	}

	if _, err := br.Exec(); err != nil {
		return canvas, fmt.Errorf("failed to insert access rule: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return canvas, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return canvas, nil
}

func (q *Queries) GetCanvasOwner(canvasID string) (string, error) {
	query := `SELECT owner_id FROM canvases WHERE canvas_id = $1`

	var ownerID string
	err := q.pool.QueryRow(context.Background(), query, canvasID).Scan(&ownerID)
	return ownerID, err
}

func (q *Queries) GetCanvas(canvasID string) (types.Canvas, error) {
	var canvas types.Canvas
	query := `SELECT * FROM canvases WHERE canvas_id = $1`

	err := q.pool.QueryRow(context.Background(), query, canvasID).Scan(
		&canvas.ID,
		&canvas.OwnerID,
		&canvas.Title,
		&canvas.Description,
		&canvas.Width,
		&canvas.Height,
		&canvas.PixelData,
		&canvas.LastEditedAt,
		&canvas.LinkAccessType,
		&canvas.LinkAccessRole,
		&canvas.CreatedAt,
		&canvas.StarCount,
	)

	return canvas, err
}

func (q *Queries) DeleteCanvas(canvasID string) error {
	query := `DELETE FROM canvases WHERE canvas_id = $1`

	_, err := q.pool.Exec(context.Background(), query, canvasID)
	return err
}

func (q *Queries) UpdateLinkAccess(canvasID string, accessType types.AccessType, accessRole types.AccessRole) error {
	if accessRole == types.Owner {
		return fmt.Errorf("access role cannot be of type owner")
	}

	query := `UPDATE canvases SET link_access_type = $1, link_access_role = $2 WHERE canvas_id = $3`
	_, err := q.pool.Exec(context.Background(), query, accessType, accessRole, canvasID)
	return err
}

func (q *Queries) GetCanvasLinkAccess(canvasID string) (types.AccessType, types.AccessRole, error) {
	query := `SELECT link_access_type, link_access_role FROM canvases WHERE canvas_id = $1`

	var accessType types.AccessType
	var accessRole types.AccessRole
	err := q.pool.QueryRow(context.Background(), query, canvasID).Scan(&accessType, &accessRole)
	return accessType, accessRole, err
}

func (q *Queries) UpdateCanvas(canvasID, title, desctription string) error {
	var values []any
	var updates []string
	argCount := 1

	if title != "" {
		values = append(values, title)
		updates = append(updates, fmt.Sprintf("title = $%d", argCount))
		argCount++
	}

	if desctription != "" {
		values = append(values, desctription)
		updates = append(updates, fmt.Sprintf("description = $%d", argCount))
		argCount++
	}

	values = append(values, canvasID)

	if len(updates) == 0 {
		return nil
	}

	query := fmt.Sprintf("UPDATE canvases SET %s WHERE canvas_id = $%d", strings.Join(updates, ","), argCount)
	_, err := q.pool.Exec(context.Background(), query, values...)
	return err
}
