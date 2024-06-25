-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name, apikey)
VALUES ($1, $2, $3, $4, encode(sha256(random()::text::bytea), 'hex'))
RETURNING *;

-- name: GetUser :one
SELECT *
FROM users
WHERE apikey = $1;

-- name: IsUser :one
SELECT EXISTS(
    SELECT 1
    FROM users
    WHERE apikey = $1
) AS user_exists;   