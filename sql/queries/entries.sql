-- ----------------
--    IP ENTRIES
-- ----------------

-- name: ListIPEntries :many
SELECT * FROM ip_entries WHERE feed_id = ?;

-- name: ListIPEntriesWindow :many
SELECT * FROM ip_entries WHERE feed_id = ? LIMIT ? OFFSET ?;

-- name: GetIPEnabledEntries :many
SELECT value FROM ip_entries WHERE enabled = ? AND feed_id = ? AND (valid_until IS NULL OR valid_until >= ?);

-- name: InsertIPEntry :execresult
INSERT INTO ip_entries (value, comment, valid_until, feed_id) VALUES (?, ?, ?, ?);

-- name: EditIPEntryById :execresult
UPDATE ip_entries SET comment = ?, valid_until = ? WHERE id = ?;

-- name: RemoveIPEntry :execresult
DELETE FROM ip_entries WHERE id = ?;

-- --------------------
--    DOMAIN ENTRIES
-- --------------------

-- name: ListDomainEntries :many
SELECT * FROM domain_entries WHERE feed_id = ?;

-- name: ListDomainEntriesWindow :many
SELECT * FROM domain_entries WHERE feed_id = ? LIMIT ? OFFSET ?;

-- name: GetDomainEnabledEntries :many
SELECT value FROM domain_entries WHERE enabled = ? AND feed_id = ? AND (valid_until IS NULL OR valid_until >= ?);

-- name: InsertDomainEntry :execresult
INSERT INTO domain_entries (value, comment, valid_until, feed_id) VALUES (?, ?, ?, ?);

-- name: EditDomainEntryById :execresult
UPDATE domain_entries SET comment = ?, valid_until = ? WHERE id = ?;

-- name: RemoveDomainEntry :execresult
DELETE FROM domain_entries WHERE id = ?;

-- -----------------
--    URL ENTRIES
-- -----------------

-- name: ListURLEntries :many
SELECT * FROM url_entries WHERE feed_id = ?;

-- name: ListURLEntriesWindow :many
SELECT * FROM url_entries WHERE feed_id = ? LIMIT ? OFFSET ?;

-- name: GetURLEnabledEntries :many
SELECT value FROM url_entries WHERE enabled = ? AND feed_id = ? AND (valid_until IS NULL OR valid_until >= ?);

-- name: InsertURLEntry :execresult
INSERT INTO url_entries (value, comment, valid_until, feed_id) VALUES (?, ?, ?, ?);

-- name: EditURLEntryById :execresult
UPDATE url_entries SET comment = ?, valid_until = ? WHERE id = ?;

-- name: RemoveURLEntry :execresult
DELETE FROM url_entries WHERE id = ?;