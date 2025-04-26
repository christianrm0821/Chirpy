-- name: GetChirpWithID :one
select * from chirps
where id = $1;