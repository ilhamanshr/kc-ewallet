-- name: CreateUser :one
INSERT INTO users (username, password, created_at)
VALUES ($1, $2, NOW())
RETURNING id;

-- name: GetUserByIDLock :one
SELECT *
FROM users
WHERE id = $1
FOR UPDATE;

-- name: GetUserByUsername :one
SELECT *
FROM users
WHERE username = $1;

-- name: UpdateUserBalanceByID :exec
UPDATE users
SET balance = $2
WHERE id = $1;