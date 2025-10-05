-- name: CreateTransaction :one
INSERT INTO transactions (user_id, amount, type, created_at)
VALUES ($1, $2, $3, NOW())
RETURNING id;