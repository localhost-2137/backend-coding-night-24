CREATE TABLE reports
(
    id         INTEGER PRIMARY KEY,
    label      TEXT NOT NULL,
    content    TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);