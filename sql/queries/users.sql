-- name: CreateNewUser :one
insert into users(id, created_at, updated_at, email, password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;


-- name: GetUserFromEmail :one
Select * from users
where email = $1;

-- name: DeleteUsers :exec
Delete from users;