ALTER TABLE posts ADD COLUMN group_id INTEGER REFERENCES user_groups(id) ON DELETE CASCADE;
CREATE INDEX IF NOT EXISTS idx_posts_group_id ON posts(group_id);
