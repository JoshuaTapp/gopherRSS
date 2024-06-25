-- name: CreatePost :one
INSERT INTO posts(id, created_at, updated_at, url, feed_id, title, description, published_at)
VALUES($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetUsersPosts :many
SELECT p.id, p.created_at, p.updated_at, p.title, p.url, p.description, p.published_at, p.feed_id
FROM posts p
INNER JOIN feed_follows f 
    ON p.feed_id = f.feed_id
WHERE f.user_id = $1
ORDER BY p.created_at DESC
LIMIT $2;   
