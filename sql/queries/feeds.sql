-- name: CreateFeed :one
INSERT INTO feeds(id, name, URL, created_at, updated_at, user_id)
VALUES($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetAllFeeds :many
SELECT *
FROM feeds
ORDER BY created_at DESC;

-- name: GetNextFeedsToFetch :many
SELECT *
FROM feeds
ORDER BY last_fetched_at IS NULL DESC, last_fetched_at ASC
LIMIT $1;

-- name: UpdateFeedFetchTime :exec
UPDATE feeds
SET last_fetched_at = NOW(), updated_at = NOW()
WHERE URL = $1;