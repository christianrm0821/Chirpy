-- name: RevokeRefreshToken :exec
update refresh_tokens
set revoked_at = $1, updated_at = $2
where token = $3;