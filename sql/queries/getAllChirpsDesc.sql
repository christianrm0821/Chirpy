-- name: GetAllChirpsDesc :many
select * from chirps
order by created_at desc;