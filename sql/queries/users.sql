-- name: ListUsers :many
SELECT * FROM users;

-- name: GetUserById :one
SELECT * FROM users WHERE id = ?;

-- name: GetUserByName :one
SELECT * FROM users WHERE name = ?;

