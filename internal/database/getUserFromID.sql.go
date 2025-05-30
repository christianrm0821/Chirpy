// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: getUserFromID.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const getUserFromID = `-- name: GetUserFromID :one
select id, created_at, updated_at, email, hashed_password, is_chirpy_red from users
where id = $1
`

func (q *Queries) GetUserFromID(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserFromID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IsChirpyRed,
	)
	return i, err
}
