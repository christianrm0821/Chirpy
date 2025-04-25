-- name: GetAllChirps :many
select * from chirps
order by created_at asc;