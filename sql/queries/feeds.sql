-- name: CreateFeed :exec
INSERT INTO feeds (
    name,
    url,
    created_at,
    updated_at
) VALUES (
    ?,
    ?,
    ?,
    ?
) RETURNING *;

-- name: GetFeedById :one
SELECT * FROM feeds WHERE id = ?;

-- name: GetFeedByUrl :one
SELECT * FROM feeds WHERE url = ?;

-- name: GetFeeds :many
SELECT * FROM feeds ORDER BY created_at ASC;

-- name: DeleteFeed :exec
DELETE FROM feeds WHERE url = ?;
