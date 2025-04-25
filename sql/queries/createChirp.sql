-- name: CreateChirp :one
Insert into chirps(id,created_at,updated_at,body,user_id)
values(
    gen_random_uuid(),
    current_timestamp,
    current_timestamp,
    $1,
    $2
)
returning *;