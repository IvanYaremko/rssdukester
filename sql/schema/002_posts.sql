-- +goose Up
CREATE TABLE posts (
id INTEGER PRIMARY KEY AUTOINCREMENT,
title TEXT NOT NULL,
url TEXT NOT NULL,
saved BOOLEAN,
viewed DATETIME NOT NULL
);
