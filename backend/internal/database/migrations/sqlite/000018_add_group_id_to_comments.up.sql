ALTER TABLE comments ADD COLUMN group_id INTEGER REFERENCES user_groups(id) ON DELETE CASCADE;
CREATE INDEX IF NOT EXISTS idx_comments_group_id ON comments(group_id);
