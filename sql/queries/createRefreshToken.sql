-- name: CreateRefreshToken :one
Insert into refresh_tokens(token, created_at,updated_at, user_id, expires_at, revoked_at)
values(
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
returning *;