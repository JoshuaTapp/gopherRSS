-- name: FollowFeed :one
INSERT INTO feed_follows(feed_id, user_id, created_at, updated_at)
VALUES($1, $2, $3, $4)
RETURNING *;
--

-- name: UnfollowFeed :exec
DELETE FROM feed_follows
WHERE 
    feed_id = $1
    AND 
    user_id = $2;
--

-- name: FollowExists :one
SELECT * FROM feed_follows WHERE feed_id = $1 AND user_id = $2;
---

-- name: GetUserFeeds :many
SELECT * 
FROM feed_follows
WHERE user_id = $1
ORDER BY updated_at DESC;
--