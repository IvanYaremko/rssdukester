-- name: CreatePost :exec
INSERT INTO posts (
    feed_id,
    title,
    url,
    content,
    published_at,
    created_at
) VALUES (
    ?,
    ?,
    ?,
    ?,
    ?,
    ?
) RETURNING *;

-- name: GetPostById :one
SELECT * FROM posts WHERE id = ?;

-- name: GetPostByUrl :one
SELECT * FROM posts WHERE url = ?;

-- name: GetPostsByFeedId :many
SELECT * FROM posts 
WHERE feed_id = ? 
ORDER BY published_at DESC;

-- name: GetRecentPosts :many
SELECT p.*, f.name as feed_name, f.url as feed_url
FROM posts p
JOIN feeds f ON p.feed_id = f.id
ORDER BY published_at DESC
LIMIT ?;
