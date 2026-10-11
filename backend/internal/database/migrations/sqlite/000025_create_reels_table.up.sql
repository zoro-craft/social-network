CREATE TABLE reels (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id    INTEGER NOT NULL,
    video      TEXT    NOT NULL,
    caption    TEXT    NOT NULL DEFAULT '',
    created_at INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE reel_likes (
    reel_id    INTEGER NOT NULL,
    user_id    INTEGER NOT NULL,
    created_at INTEGER NOT NULL,
    PRIMARY KEY (reel_id, user_id),
    FOREIGN KEY (reel_id) REFERENCES reels(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE reel_views (
    reel_id   INTEGER NOT NULL,
    user_id   INTEGER NOT NULL,
    viewed_at INTEGER NOT NULL,
    PRIMARY KEY (reel_id, user_id),
    FOREIGN KEY (reel_id) REFERENCES reels(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- A save is a per-user bookmark, so one row per (reel, user).
CREATE TABLE reel_saves (
    reel_id    INTEGER NOT NULL,
    user_id    INTEGER NOT NULL,
    created_at INTEGER NOT NULL,
    PRIMARY KEY (reel_id, user_id),
    FOREIGN KEY (reel_id) REFERENCES reels(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Shares accumulate: every share is its own row, so the count reflects events
-- rather than unique sharers.
CREATE TABLE reel_shares (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    reel_id    INTEGER NOT NULL,
    user_id    INTEGER NOT NULL,
    created_at INTEGER NOT NULL,
    FOREIGN KEY (reel_id) REFERENCES reels(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE reel_comments (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    reel_id    INTEGER NOT NULL,
    user_id    INTEGER NOT NULL,
    content    TEXT    NOT NULL,
    created_at INTEGER NOT NULL,
    FOREIGN KEY (reel_id) REFERENCES reels(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_reels_created ON reels (created_at DESC);
CREATE INDEX idx_reels_user_created ON reels (user_id, created_at DESC);
CREATE INDEX idx_reel_comments_reel ON reel_comments (reel_id, created_at DESC);
CREATE INDEX idx_reel_shares_reel ON reel_shares (reel_id);
CREATE INDEX idx_reel_saves_user ON reel_saves (user_id, created_at DESC);
