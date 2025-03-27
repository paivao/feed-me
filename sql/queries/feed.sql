-- name: ListFeeds :many
SELECT * FROM feeds;

-- name: GetFeedById :one
SELECT * FROM feeds WHERE id = $1;

-- name: GetFeedByName :one
SELECT * FROM feeds WHERE name = $1;

-- name: GetFeedByIdAndType :one
SELECT * FROM feeds WHERE id = $1 and type = $2;

-- name: GetFeedTypeById :one
SELECT type FROM feeds WHERE id = $1;

-- name: CreateFeed :one
INSERT INTO feeds (name, description, is_public, type) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: EditFeedById :exec
UPDATE feeds SET description = $1, is_public = $2 WHERE id = $3;

-- name: RemoveFeedById :exec
DELETE FROM feeds WHERE id = $1;