CREATE TABLE notifications_new (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  actor_id INTEGER,
  type TEXT NOT NULL,
  ref_id INTEGER,
  message TEXT NOT NULL,
  is_read BOOLEAN NOT NULL DEFAULT 0,
  created_at INTEGER,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (actor_id) REFERENCES users(id) ON DELETE SET NULL
);

INSERT INTO notifications_new (id, user_id, actor_id, type, ref_id, message, is_read, created_at)
  SELECT id, user_id, NULL,
         CASE WHEN type = 'group_invtation' THEN 'group_invitation' ELSE type END,
         COALESCE(request_id, invitation_id), message, is_read,
         CAST(strftime('%s', created_at) AS INTEGER) * 1000
  FROM notifications;

DROP TABLE notifications;
ALTER TABLE notifications_new RENAME TO notifications;