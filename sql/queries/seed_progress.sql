-- name: GetSeedProgress :one
SELECT id, last_processed_index
FROM seed_progress
LIMIT 1;

-- name: UpdateSeedProgress :exec
UPDATE seed_progress
SET last_processed_index = :last_processed_index,
    updated_at = DATETIME('now', 'utc')
WHERE id = :id;