-- +goose Up
CREATE TABLE saved_posts (
id INTEGER PRIMARY KEY AUTOINCREMENT,
title TEXT NOT NULL,
url TEXT NOT NULL,
feed TEXT NOT NULL,
created_at DATETIME NOT NULL,
UNIQUE(url)
);

-- +goose Down
DROP TABLE saved_posts;
