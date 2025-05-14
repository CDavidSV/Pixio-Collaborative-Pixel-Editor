package data

import (
	"context"

	"github.com/CDavidSV/Pixio/types"
)

func (q *Queries) GetUserAccess(objectID string, objectType types.ObjectType, userID string) (types.UserAccess, error) {
	query := `SELECT * FROM user_access WHERE object_id = $1 AND object_type = $2 AND user_id = $3`

	var userAccess types.UserAccess
	err := q.pool.QueryRow(context.Background(), query, objectID, objectType, userID).Scan(
		&userAccess.ObjectID,
		&userAccess.ObjectType,
		&userAccess.UserID,
		&userAccess.AccessRole,
		&userAccess.LastModifiedAt,
		&userAccess.LastModifiedBy,
	)

	return userAccess, err
}

func (q *Queries) CreateUserAccess(objectID string, objectType types.ObjectType, accessRole types.AccessRole, userID, creatorUserID string) (types.UserAccess, error) {
	query := `INSERT INTO user_access (object_id, object_type, user_id, access_role, last_modified_by) VALUES ($1, $2, $3, $4, $5) RETURNING last_modified_at`

	userAccess := types.UserAccess{
		ObjectID:       objectID,
		ObjectType:     objectType,
		UserID:         userID,
		AccessRole:     accessRole,
		LastModifiedBy: creatorUserID,
	}
	err := q.pool.QueryRow(context.Background(), query, objectID, objectType, userID, accessRole, creatorUserID).Scan(&userAccess.LastModifiedAt)
	return userAccess, err
}
