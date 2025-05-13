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
