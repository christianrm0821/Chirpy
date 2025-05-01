-- name: DeleteChirpWithID :exec
delete from chirps
where id =$1;