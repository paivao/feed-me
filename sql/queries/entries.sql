-- ----------------
--    IP ENTRIES
-- ----------------

-- name: ListIPEntries :many
SELECT * FROM ip_entries WHERE feed_id = ?;

-- name: GetIPEntryById :one
SELECT * FROM ip_entries WHERE id = ? AND feed_id = ?;

-- name: ListIPEntriesWindow :many
SELECT * FROM ip_entries WHERE feed_id = ? LIMIT ? OFFSET ?;

-- name: GetIPEnabledEntries :many
SELECT value FROM ip_entries WHERE enabled = 1 AND feed_id = ? AND (valid_until IS NULL OR valid_until >= ?);

-- name: InsertIPEntry :execresult
INSERT INTO ip_entries (value, description, valid_until, feed_id) VALUES (?, ?, ?, ?);

-- name: EditIPEntryById :exec
UPDATE ip_entries SET enabled = ?, description = ?, valid_until = ? WHERE id = ?;

-- name: RemoveIPEntry :exec
DELETE FROM ip_entries WHERE id = ? AND feed_id = ?;

-- --------------------
--    DOMAIN ENTRIES
-- --------------------

-- name: ListDomainEntries :many
SELECT * FROM domain_entries WHERE feed_id = ?;

-- name: GetDomainEntryById :one
SELECT * FROM domain_entries WHERE id = ? AND feed_id = ?;

-- name: ListDomainEntriesWindow :many
SELECT * FROM domain_entries WHERE feed_id = ? LIMIT ? OFFSET ?;

-- name: GetDomainEnabledEntries :many
SELECT value FROM domain_entries WHERE enabled = 1 AND feed_id = ? AND (valid_until IS NULL OR valid_until >= ?);

-- name: InsertDomainEntry :execresult
INSERT INTO domain_entries (value, description, valid_until, feed_id) VALUES (?, ?, ?, ?);

-- name: EditDomainEntryById :exec
UPDATE domain_entries SET enabled = ?, description = ?, valid_until = ? WHERE id = ?;

-- name: RemoveDomainEntry :exec
DELETE FROM domain_entries WHERE id = ? AND feed_id = ?;

-- -----------------
--    URL ENTRIES
-- -----------------

-- name: ListURLEntries :many
SELECT * FROM url_entries WHERE feed_id = ?;

-- name: GetURLEntryById :one
SELECT * FROM url_entries WHERE id = ? AND feed_id = ?;

-- name: ListURLEntriesWindow :many
SELECT * FROM url_entries WHERE feed_id = ? LIMIT ? OFFSET ?;

-- name: GetURLEnabledEntries :many
SELECT value FROM url_entries WHERE enabled = 1 AND feed_id = ? AND (valid_until IS NULL OR valid_until >= ?);

-- name: InsertURLEntry :execresult
INSERT INTO url_entries (value, description, valid_until, feed_id) VALUES (?, ?, ?, ?);

-- name: EditURLEntryById :exec
UPDATE url_entries SET enabled = ?, description = ?, valid_until = ? WHERE id = ?;

-- name: RemoveURLEntry :exec
DELETE FROM url_entries WHERE id = ? AND feed_id = ?;