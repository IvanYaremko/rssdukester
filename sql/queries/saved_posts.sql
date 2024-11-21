-- name: SavePost :exec
INSERT INTO saved_posts (
title,
url,
feed,
created_at
) VALUES (
?,
?,
?,
?
);

-- name: GetSavedPosts :many
SELECT * FROM saved_posts ORDER BY created_at DESC;