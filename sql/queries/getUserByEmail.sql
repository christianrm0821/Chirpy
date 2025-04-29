-- name: GetUserByEmail :one
select * from users
where email = $1;