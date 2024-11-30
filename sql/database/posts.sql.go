// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: posts.sql

package database

import (
	"context"
	"database/sql"
	"time"
)

const createPost = `-- name: CreatePost :exec
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
)
`

type CreatePostParams struct {
	FeedID      int64
	Title       string
	Url         string
	Content     sql.NullString
	PublishedAt time.Time
	LastViewed  time.Time
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) error {
	_, err := q.db.ExecContext(ctx, createPost,
		arg.FeedID,
		arg.Title,
		arg.Url,
		arg.Content,
		arg.PublishedAt,
		arg.LastViewed,
	)
	return err
}

const deletePost = `-- name: DeletePost :exec
DELETE FROM posts WHERE url = ?
`

func (q *Queries) DeletePost(ctx context.Context, url string) error {
	_, err := q.db.ExecContext(ctx, deletePost, url)
	return err
}

const getPostById = `-- name: GetPostById :one
SELECT id, feed_id, title, url, content, published_at, last_viewed FROM posts WHERE id = ?
`

func (q *Queries) GetPostById(ctx context.Context, id int64) (Post, error) {
	row := q.db.QueryRowContext(ctx, getPostById, id)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.FeedID,
		&i.Title,
		&i.Url,
		&i.Content,
		&i.PublishedAt,
		&i.LastViewed,
	)
	return i, err
}

const getPostByUrl = `-- name: GetPostByUrl :one
SELECT id, feed_id, title, url, content, published_at, last_viewed FROM posts WHERE url = ?
`

func (q *Queries) GetPostByUrl(ctx context.Context, url string) (Post, error) {
	row := q.db.QueryRowContext(ctx, getPostByUrl, url)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.FeedID,
		&i.Title,
		&i.Url,
		&i.Content,
		&i.PublishedAt,
		&i.LastViewed,
	)
	return i, err
}

const getPostContent = `-- name: GetPostContent :exec
SELECT content FROM posts WHERE url = ? AND content IS NOT NULL
`

func (q *Queries) GetPostContent(ctx context.Context, url string) error {
	_, err := q.db.ExecContext(ctx, getPostContent, url)
	return err
}

const getPosts = `-- name: GetPosts :many
SELECT id, feed_id, title, url, content, published_at, last_viewed FROM posts ORDER BY published_at DESC
`

func (q *Queries) GetPosts(ctx context.Context) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, getPosts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.FeedID,
			&i.Title,
			&i.Url,
			&i.Content,
			&i.PublishedAt,
			&i.LastViewed,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPostsByFeedId = `-- name: GetPostsByFeedId :many
SELECT id, feed_id, title, url, content, published_at, last_viewed FROM posts
WHERE feed_id = ?
ORDER BY published_at DESC
`

func (q *Queries) GetPostsByFeedId(ctx context.Context, feedID int64) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, getPostsByFeedId, feedID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.FeedID,
			&i.Title,
			&i.Url,
			&i.Content,
			&i.PublishedAt,
			&i.LastViewed,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPostsWithFeed = `-- name: GetPostsWithFeed :many
SELECT p.id, p.feed_id, p.title, p.url, p.content, p.published_at, p.last_viewed, f.name as feed_name, f.url as feed_url
FROM posts p
JOIN feeds f ON p.feed_id = f.id
ORDER BY published_at DESC
`

type GetPostsWithFeedRow struct {
	ID          int64
	FeedID      int64
	Title       string
	Url         string
	Content     sql.NullString
	PublishedAt time.Time
	LastViewed  time.Time
	FeedName    string
	FeedUrl     string
}

func (q *Queries) GetPostsWithFeed(ctx context.Context) ([]GetPostsWithFeedRow, error) {
	rows, err := q.db.QueryContext(ctx, getPostsWithFeed)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPostsWithFeedRow
	for rows.Next() {
		var i GetPostsWithFeedRow
		if err := rows.Scan(
			&i.ID,
			&i.FeedID,
			&i.Title,
			&i.Url,
			&i.Content,
			&i.PublishedAt,
			&i.LastViewed,
			&i.FeedName,
			&i.FeedUrl,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateLastViewed = `-- name: UpdateLastViewed :exec
UPDATE posts SET last_viewed = ? WHERE url = ?
`

type UpdateLastViewedParams struct {
	LastViewed time.Time
	Url        string
}

func (q *Queries) UpdateLastViewed(ctx context.Context, arg UpdateLastViewedParams) error {
	_, err := q.db.ExecContext(ctx, updateLastViewed, arg.LastViewed, arg.Url)
	return err
}

const updatePostContent = `-- name: UpdatePostContent :exec
UPDATE posts SET content = ? WHERE url = ?
`

type UpdatePostContentParams struct {
	Content sql.NullString
	Url     string
}

func (q *Queries) UpdatePostContent(ctx context.Context, arg UpdatePostContentParams) error {
	_, err := q.db.ExecContext(ctx, updatePostContent, arg.Content, arg.Url)
	return err
}
