DROP INDEX IF EXISTS idx_comments_group_id;
ALTER TABLE comments DROP COLUMN group_id;