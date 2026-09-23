CREATE TABLE messages (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    sender_id  INTEGER NOT NULL,
    type       TEXT    NOT NULL CHECK (type IN ('private','group')),
    ref_id     INTEGER NOT NULL,
    content    TEXT    NOT NULL,
    created_at INTEGER NOT NULL,
    FOREIGN KEY (sender_id) REFERENCES users(id)
);