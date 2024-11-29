-- name: CreateSavedPost :exec
INSERT INTO saved_posts (
    post_id,
    created_at
) VALUES (
    ?,
    ?
);

-- name: GetSavedPosts :many
SELECT
    p.*,
    f.name as feed_name,
    f.url as feed_url,
    sp.created_at as saved_at
FROM saved_posts sp
JOIN posts p ON sp.post_id = p.id
JOIN feeds f ON p.feed_id = f.id
ORDER BY sp.created_at DESC;

-- name: GetSavedPost :one
SELECT
    p.*,
    f.name as feed_name,
    f.url as feed_url,
    sp.created_at as saved_at
FROM saved_posts sp
JOIN posts p ON sp.post_id = p.id
JOIN feeds f ON p.feed_id = f.id
WHERE p.url = ?;

-- name: DeleteSavedPost :exec
DELETE FROM saved_posts WHERE post_id = (
    SELECT id FROM posts WHERE url = ?
);

-- name: IsPostSaved :one
SELECT EXISTS(
    SELECT 1 FROM saved_posts sp
    JOIN posts p ON sp.post_id = p.id
    WHERE p.url = ?
);
