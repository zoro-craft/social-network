package posts

import "zone/internal/database"


const visibleCondition = `(
	p.privacy = 'public'
	OR p.user_id = ?
	OR (p.privacy = 'almost_private' AND EXISTS (
		SELECT 1 FROM followers f
		WHERE f.follower_id = ? AND f.followee_id = p.user_id AND f.status = 'accepted'
	))
	OR (p.privacy = 'private' AND EXISTS (
		SELECT 1 FROM post_viewers pv WHERE pv.post_id = p.id AND pv.user_id = ?
	))
)`

func CanViewPost(postID, viewerID int) bool {
	var exists bool
	database.Database.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM posts p WHERE p.id = ? AND `+visibleCondition+`)`,
		postID, viewerID, viewerID, viewerID,
	).Scan(&exists)
	return exists
}