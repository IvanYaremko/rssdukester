-- name: CreateFeed :exec
INSERT INTO feeds (
    id,
    name,
    url,
    created_at,
    updated_at
) VALUES (
    ?,
    ?,
    ?,
    ?,
    ?
) RETURNING *;

-- name: GetFeedById :one
SELECT * FROM feeds WHERE id = ?;

-- name: GetFeedByUrl :one
SELECT * FROM feeds WHERE url = ?;

