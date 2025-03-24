-- name: ListFeeds :many
SELECT * FROM feeds;

-- name: GetFeedById :one
SELECT * FROM feeds WHERE id = ?;

-- name: GetFeedByName :one
SELECT * FROM feeds WHERE name = ?;

-- name: GetFeedByNameAndType :one
SELECT * FROM feeds WHERE name = ? and type = ?;

-- name: GetFeedTypeById :one
SELECT type FROM feeds WHERE id = ?;

-- name: CreateFeed :execresult
INSERT INTO feeds (name, comment, is_public, type) VALUES (?, ?, ?, ?);

-- name: EditFeedById :execresult
UPDATE feeds SET comment = ?, is_public = ? WHERE id = ?;

-- name: RemoveFeedById :execresult
DELETE FROM feeds WHERE id = ?;