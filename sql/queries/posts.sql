-- name: CreatePost :exec
INSERT INTO posts (
feed_id,
title,
url,
content,
published_at,
last_viewed
) VALUES (
?,
?,
?,
?,
?,
?
);

-- name: GetPostById :one
SELECT * FROM posts WHERE id = ?;

-- name: GetPostByUrl :one
SELECT * FROM posts WHERE url = ?;

-- name: GetPostsByFeedId :many
SELECT * FROM posts
WHERE feed_id = ?
ORDER BY published_at DESC;

-- name: GetPosts :many
SELECT * FROM posts ORDER BY published_at DESC;

-- name: GetPostsWithFeed :many
SELECT p.*, f.name as feed_name, f.url as feed_url
FROM posts p
JOIN feeds f ON p.feed_id = f.id
ORDER BY published_at DESC;

-- name: UpdatePostContent :exec
UPDATE posts SET content = ? WHERE url = ?;

-- name: UpdateLastViewed :exec
UPDATE posts SET last_viewed = ? WHERE url = ?;

-- name: GetPostContent :exec
SELECT content FROM posts WHERE url = ? AND content IS NOT NULL;

-- name: DeletePost :exec
DELETE FROM posts WHERE url = ?;
