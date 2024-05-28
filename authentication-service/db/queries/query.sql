-- name: CreateUser :one
INSERT INTO users (
    email, first_name, last_name, password, active
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY first_name;

-- name: UpdateUser :exec
UPDATE users
    set first_name = $1,
    last_name = $2,
    password = $3,
    active = $4,
    update_at = $5
WHERE id = $6;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
