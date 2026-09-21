package repositories

import "zone/internal/database"

// GetPresenceAudience returns the userIDs entitled to see the user's presence.
func GetPresenceAudience(userID int) []int {
	rows, err := database.Database.Query(`
		SELECT follower_id
		FROM followers
		WHERE followee_id = ? AND status = 'accepted'
		UNION
		SELECT followee_id
		FROM followers
		WHERE follower_id = ? AND status = 'accepted'
	`, userID, userID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil
	}
	return ids
}
