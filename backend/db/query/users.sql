-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: CreateUser :one
INSERT INTO users (name, email, age) VALUES ($1, $2, $3) RETURNING *;

-- name: DeleteUserByID :one
DELETE FROM users WHERE id = $1 RETURNING *;

-- name: DeleteUserByEmail :one
DELETE FROM users WHERE email = $1 RETURNING *;

-- name: UpdateUser :one
UPDATE users SET name = $2, email = $3, age = $4 WHERE id = $1 RETURNING *;
