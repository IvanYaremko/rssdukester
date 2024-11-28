-- +goose Up
CREATE TABLE saved_posts (
id INTEGER PRIMARY KEY AUTOINCREMENT,
post_id INTEGER NOT NULL,
created_at DATETIME NOT NULL,
UNIQUE(post_id),
FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE saved_posts;
