CREATE TABLE notifications_old (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  type TEXT NOT NULL,
  request_id INTEGER,
  invitation_id INTEGER,
  message TEXT NOT NULL,
  is_read BOOLEAN NOT NULL DEFAULT 0,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (invitation_id) REFERENCES group_invitations(id) ON DELETE CASCADE
);

INSERT INTO notifications_old (id, user_id, type, request_id, invitation_id, message, is_read, created_at)
  SELECT id, user_id,
         CASE WHEN type = 'group_invitation' THEN 'group_invtation' ELSE type END,
         CASE WHEN type = 'join_request' THEN ref_id ELSE NULL END,
         CASE WHEN type = 'group_invitation' THEN ref_id ELSE NULL END,
         message, is_read,
         datetime(created_at / 1000, 'unixepoch')
  FROM notifications;

DROP TABLE notifications;
ALTER TABLE notifications_old RENAME TO notifications;