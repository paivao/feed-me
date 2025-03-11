-- name: ListFeeds :many
SELECT * FROM feeds;

-- name: GetFeedById :one
SELECT * FROM feeds WHERE id = ?;

-- name: GetFeedByName :one
SELECT * FROM feeds WHERE name = ?;

-- name: GetFeedTypeById :one
SELECT type FROM feeds WHERE id = ?;

-- name: ListIPEntries :many
SELECT * FROM ip_entries WHERE feed_id = ?;

-- name: GetIPEnabledEntries :many
SELECT value FROM ip_entries WHERE enabled = ? AND feed_id = ? AND (valid_until IS NULL OR valid_until >= ?);

-- name: GetDomainEnabledEntries :many
SELECT value FROM domain_entries WHERE enabled = ? AND feed_id = ? AND (valid_until IS NULL OR valid_until >= ?);

-- name: GetURLEnabledEntries :many
SELECT value FROM url_entries WHERE enabled = ? AND feed_id = ? AND (valid_until IS NULL OR valid_until >= ?);

-- name: ListUsers :many
SELECT * FROM users;

-- name: GetUserById :one
SELECT * FROM users WHERE id = ?;

-- name: GetUserByName :one
SELECT * FROM users WHERE name = ?;

-- name: CreateFeed :exec
INSERT INTO feeds (name, is_public, type) VALUES (?, ?, ?);

-- name: InsertIPEntry :exec
INSERT INTO ip_entries (value, valid_until, feed_id) VALUES (?, ?, ?);