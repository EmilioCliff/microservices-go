-- name: CreateLog :one
INSERT INTO logs (
    email, data, user_ip, user_agent
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetLogByID :one
SELECT * FROM logs
WHERE id = $1 LIMIT 1;

-- name: GetLogByEmail :many
SELECT * FROM logs
WHERE email = $1
ORDER BY logged_at DESC;

-- name: GetLogByIp :many
SELECT * FROM logs
WHERE user_ip = $1
ORDER BY logged_at DESC;

-- name: ListLogs :many
SELECT * FROM logs
ORDER BY logged_at DESC;

