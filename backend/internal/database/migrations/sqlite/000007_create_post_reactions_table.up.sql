CREATE TABLE IF NOT EXISTS post_reactions (
    user_id INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    is_like INTEGER NOT NULL DEFAULT 1,
    PRIMARY KEY (user_id, post_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (post_id) REFERENCES posts(id)
);