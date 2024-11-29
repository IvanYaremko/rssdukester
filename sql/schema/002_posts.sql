-- +goose Up
CREATE TABLE posts (
id INTEGER PRIMARY KEY AUTOINCREMENT,
feed_id INTEGER NOT NULL, 
title TEXT NOT NULL,
url TEXT NOT NULL,
content TEXT,
published_at DATETIME NOT NULL,
last_viewed DATETIME NOT NULL,
UNIQUE(url),
FOREIGN KEY(feed_id) REFERENCES feeds(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE posts;
