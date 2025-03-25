-- name: ListFeeds :many
SELECT * FROM feeds;

-- name: GetFeedById :one
SELECT * FROM feeds WHERE id = ?;

-- name: GetFeedByName :one
SELECT * FROM feeds WHERE name = ?;

-- name: GetFeedByIdAndType :one
SELECT * FROM feeds WHERE id = ? and type = ?;

-- name: GetFeedTypeById :one
SELECT type FROM feeds WHERE id = ?;

-- name: CreateFeed :execresult
INSERT INTO feeds (name, description, is_public, type) VALUES (?, ?, ?, ?);

-- name: EditFeedById :exec
UPDATE feeds SET description = ?, is_public = ? WHERE id = ?;

-- name: RemoveFeedById :exec
DELETE FROM feeds WHERE id = ?;