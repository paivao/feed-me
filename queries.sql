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

-- name: ListDomainEntries :many
SELECT * FROM domain_entries WHERE feed_id = ?;

-- name: ListURLEntries :many
SELECT * FROM url_entries WHERE feed_id = ?;

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

-- name: CreateFeed :execresult
INSERT INTO feeds (name, comment, is_public, type) VALUES (?, ?, ?, ?);

-- name: EditFeedById :execresult
UPDATE feeds SET comment = ?, is_public = ? WHERE id = ?;

-- name: RemoveFeedById :execresult
DELETE FROM feeds WHERE id = ?;

-- name: InsertIPEntry :execresult
INSERT INTO ip_entries (value, comment, valid_until, feed_id) VALUES (?, ?, ?, ?);

-- name: InsertDomainEntry :execresult
INSERT INTO domain_entries (value, comment, valid_until, feed_id) VALUES (?, ?, ?, ?);

-- name: InsertURLEntry :execresult
INSERT INTO url_entries (value, comment, valid_until, feed_id) VALUES (?, ?, ?, ?);

-- name: EditIPEntryById :execresult
UPDATE ip_entries SET comment = ?, valid_until = ? WHERE id = ?;

-- name: EditDomainEntryById :execresult
UPDATE domain_entries SET comment = ?, valid_until = ? WHERE id = ?;

-- name: EditURLEntryById :execresult
UPDATE url_entries SET comment = ?, valid_until = ? WHERE id = ?;

-- name: RemoveIPEntry :execresult
DELETE FROM ip_entries WHERE id = ?;

-- name: RemoveDomainEntry :execresult
DELETE FROM domain_entries WHERE id = ?;

-- name: RemoveURLEntry :execresult
DELETE FROM url_entries WHERE id = ?;