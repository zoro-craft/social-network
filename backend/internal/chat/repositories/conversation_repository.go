package repositories

import (
	"database/sql"

	"zone/internal/chat/models"
	"zone/internal/database"
)

// FollowGraphUsers lists every user with an accepted follow in either direction, plus every
// user with whom the caller has ever exchanged a private message (kept visible after an
// unfollow), including per-entry can_chat via the same gate as sending.
func (r ChatRepository) FollowGraphUsers(userID int) ([]models.ConversationUser, error) {
	rows, err := database.Database.Query(`
		SELECT DISTINCT u.id, u.full_name, COALESCE(u.avatar, '')
		FROM users u
		WHERE u.id IN (
			SELECT follower_id FROM followers WHERE followee_id = ? AND status = 'accepted'
			UNION
			SELECT followee_id FROM followers WHERE follower_id = ? AND status = 'accepted'
			UNION
			SELECT sender_id FROM messages WHERE type = 'private' AND ref_id = ?
			UNION
			SELECT ref_id FROM messages WHERE type = 'private' AND sender_id = ?
		)`,
		userID, userID, userID, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []models.ConversationUser{}
	for rows.Next() {
		var u models.ConversationUser
		if err := rows.Scan(&u.UserID, &u.Name, &u.Avatar); err != nil {
			return nil, err
		}
		u.CanChat = r.CanSendDM(userID, u.UserID)
		users = append(users, u)
	}
	return users, rows.Err()
}

// GroupMemberships lists every group the user belongs to, empty or not.
func (r ChatRepository) GroupMemberships(userID int) ([]models.ConversationGroup, error) {
	rows, err := database.Database.Query(`
		SELECT g.id, g.title
		FROM group_members gm
		JOIN user_groups g ON g.id = gm.group_id
		WHERE gm.user_id = ?`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := []models.ConversationGroup{}
	for rows.Next() {
		var g models.ConversationGroup
		if err := rows.Scan(&g.GroupID, &g.Title); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, rows.Err()
}
