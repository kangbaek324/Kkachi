-- name: CreateUser :one
INSERT INTO users (username, password) VALUES ($1, $2) RETURNING *;
-- name: GetUser :one
SELECT * FROM users WHERE username = $1;


-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES($1, $2, $3) RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens WHERE token = $1;