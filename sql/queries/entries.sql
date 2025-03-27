-- ----------------
--    IP ENTRIES
-- ----------------

-- name: ListIPEntries :many
SELECT * FROM ip_entries WHERE feed_id = $1;

-- name: GetIPEntryById :one
SELECT * FROM ip_entries WHERE id = $1 AND feed_id = $2;

-- name: ListIPEntriesWindow :many
SELECT * FROM ip_entries WHERE feed_id = $1 LIMIT $2 OFFSET $3;

-- name: GetIPEnabledEntries :many
SELECT value FROM ip_entries WHERE enabled = TRUE AND feed_id = $1 AND (valid_until IS NULL OR valid_until >= $2);

-- name: InsertIPEntry :one
INSERT INTO ip_entries (value, description, valid_until, feed_id) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: EditIPEntryById :exec
UPDATE ip_entries SET enabled = $1, description = $2, valid_until = $3 WHERE id = $4 RETURNING *;

-- name: RemoveIPEntry :exec
DELETE FROM ip_entries WHERE id = $1 AND feed_id = $2;

-- --------------------
--    DOMAIN ENTRIES
-- --------------------

-- name: ListDomainEntries :many
SELECT * FROM domain_entries WHERE feed_id = $1;

-- name: GetDomainEntryById :one
SELECT * FROM domain_entries WHERE id = $1 AND feed_id = $2;

-- name: ListDomainEntriesWindow :many
SELECT * FROM domain_entries WHERE feed_id = $1 LIMIT $2 OFFSET $3;

-- name: GetDomainEnabledEntries :many
SELECT value FROM domain_entries WHERE enabled = TRUE AND feed_id = $1 AND (valid_until IS NULL OR valid_until >=$2);

-- name: InsertDomainEntry :one
INSERT INTO domain_entries (value, description, valid_until, feed_id) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: EditDomainEntryById :exec
UPDATE domain_entries SET enabled = $1, description = $2, valid_until = $3 WHERE id = $4 RETURNING *;

-- name: RemoveDomainEntry :exec
DELETE FROM domain_entries WHERE id = $1 AND feed_id = $2;

-- -----------------
--    URL ENTRIES
-- -----------------

-- name: ListURLEntries :many
SELECT * FROM url_entries WHERE feed_id = $1;

-- name: GetURLEntryById :one
SELECT * FROM url_entries WHERE id = $1 AND feed_id = $2;

-- name: ListURLEntriesWindow :many
SELECT * FROM url_entries WHERE feed_id = $1 LIMIT $2 OFFSET $3;

-- name: GetURLEnabledEntries :many
SELECT value FROM url_entries WHERE enabled = TRUE AND feed_id = $1 AND (valid_until IS NULL OR valid_until >= $2);

-- name: InsertURLEntry :one
INSERT INTO url_entries (value, description, valid_until, feed_id) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: EditURLEntryById :exec
UPDATE url_entries SET enabled = $1, description = $2, valid_until = $3 WHERE id = $4 RETURNING *;

-- name: RemoveURLEntry :exec
DELETE FROM url_entries WHERE id = $1 AND feed_id = $2;