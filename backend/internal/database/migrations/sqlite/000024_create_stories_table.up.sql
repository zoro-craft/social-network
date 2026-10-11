CREATE TABLE stories (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id    INTEGER NOT NULL,
    image      TEXT    NOT NULL,
    caption    TEXT    NOT NULL DEFAULT '',
    created_at INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE story_views (
    story_id  INTEGER NOT NULL,
    user_id   INTEGER NOT NULL,
    viewed_at INTEGER NOT NULL,
    PRIMARY KEY (story_id, user_id),
    FOREIGN KEY (story_id) REFERENCES stories(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_stories_user_created ON stories (user_id, created_at DESC);
