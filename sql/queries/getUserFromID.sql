-- name: GetUserFromID :one
select * from users
where id = $1;
