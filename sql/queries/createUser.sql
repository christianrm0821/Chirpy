-- name: CreateUser :one
Insert into users(id,created_at,updated_at,email)
values(
    gen_random_uuid(),
    current_timestamp,
    current_timestamp,
    $1
)
returning *;