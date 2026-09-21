package posts

import "zone/internal/database"

const publicOrGroupCondition = `(
	(p.group_id IS NULL AND p.privacy = 'public')
	OR (p.group_id IS NOT NULL AND EXISTS (
		SELECT * FROM group_members gm WHERE gm.group_id = p.group_id AND gm.user_id = ?
	))
)`

const ownPostCondition = `p.user_id = ?`

const almostPrivateCondition = `(
	p.privacy = 'almost_private' AND EXISTS (
		SELECT * FROM followers f
		WHERE f.follower_id = ? AND f.followee_id = p.user_id AND f.status = 'accepted'
	)
)`

const privateCondition = `(
	p.privacy = 'private' AND EXISTS (
		SELECT * FROM post_viewers pv WHERE pv.post_id = p.id AND pv.user_id = ?
	)
)`

const visibleCondition = `(` +
	publicOrGroupCondition + ` OR ` +
	ownPostCondition + ` OR ` +
	almostPrivateCondition + ` OR ` +
	privateCondition +
	`)`

func CanViewPost(postID, viewerID int) bool {
	var exists bool
	database.Database.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM posts p WHERE p.id = ? AND `+visibleCondition+`)`,
		postID, viewerID, viewerID, viewerID, viewerID,
	).Scan(&exists)
	return exists
}
