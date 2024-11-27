-- +goose Up
-- First create the new posts table
CREATE TABLE posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    feed_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    content TEXT,
    published_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL,
    UNIQUE(url),
    FOREIGN KEY(feed_id) REFERENCES feeds(id) ON DELETE CASCADE
);

-- Create new saved_posts table
CREATE TABLE new_saved_posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER NOT NULL,
    created_at DATETIME NOT NULL,
    UNIQUE(post_id),
    FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE
);

-- First, let's ensure we have all the necessary feeds
INSERT OR IGNORE INTO feeds (name, url, created_at, updated_at)
SELECT DISTINCT 
    feed as name,
    feed as url,
    datetime('now') as created_at,
    datetime('now') as updated_at
FROM saved_posts
WHERE feed NOT IN (SELECT url FROM feeds);

-- Now migrate existing saved posts data
INSERT INTO posts (feed_id, title, url, content, published_at, created_at)
SELECT 
    (SELECT id FROM feeds WHERE feeds.url = saved_posts.feed),
    saved_posts.title,
    saved_posts.url,
    NULL,
    saved_posts.created_at,
    saved_posts.created_at
FROM saved_posts;

-- Then populate new_saved_posts with references to the new posts
INSERT INTO new_saved_posts (post_id, created_at)
SELECT 
    posts.id,
    saved_posts.created_at
FROM saved_posts
JOIN posts ON posts.url = saved_posts.url;

-- Drop the old saved_posts table
DROP TABLE saved_posts;

-- Rename new_saved_posts to saved_posts
ALTER TABLE new_saved_posts RENAME TO saved_posts;

-- Create the indexes
CREATE INDEX idx_posts_feed_id ON posts(feed_id);
CREATE INDEX idx_posts_published_at ON posts(published_at);

-- +goose Down
-- Remove indexes
DROP INDEX IF EXISTS idx_posts_feed_id;
DROP INDEX IF EXISTS idx_posts_published_at;

-- Recreate original saved_posts table
CREATE TABLE old_saved_posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    feed TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    UNIQUE(url)
);

-- Migrate data back to the original format
INSERT INTO old_saved_posts (title, url, feed, created_at)
SELECT 
    p.title,
    p.url,
    f.url,
    sp.created_at
FROM saved_posts sp
JOIN posts p ON sp.post_id = p.id
JOIN feeds f ON p.feed_id = f.id;

-- Drop new tables
DROP TABLE saved_posts;
DROP TABLE posts;

-- Rename old_saved_posts back to saved_posts
ALTER TABLE old_saved_posts RENAME TO saved_posts;
