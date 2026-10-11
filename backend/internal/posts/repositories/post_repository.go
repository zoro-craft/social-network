package repositories

import (
	"zone/internal/database"
	"zone/internal/posts/models"
)

const PostSelect = `
	SELECT
		p.id,
		p.title,
		p.content,
		COALESCE(p.image, ''),
		p.privacy,
		p.group_id,
		p.user_id,
		u.full_name,
		COALESCE(u.avatar, ''),
		p.created_at
`

const ReactionSelect = `,
	(SELECT COUNT(*)
	 FROM post_reactions
	 WHERE post_id = p.id AND is_like = 1) AS like_count,

	(SELECT COUNT(*)
	 FROM post_reactions
	 WHERE post_id = p.id AND is_like = 0) AS dislike_count
`

// The reacting viewer id fills the single `?` below.
const UserReactionSelect = `,
	COALESCE((
		SELECT CASE WHEN is_like = 1 THEN 'like' ELSE 'dislike' END
		FROM post_reactions
		WHERE post_id = p.id AND user_id = ?
	), '') AS user_reaction
`

const FromPosts = `
	FROM posts p
	JOIN users u ON p.user_id = u.id
`

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

const VisibleCondition = `(` +
	publicOrGroupCondition + ` OR ` +
	ownPostCondition + ` OR ` +
	almostPrivateCondition + ` OR ` +
	privateCondition +
	`)`

func CanViewPost(postID, viewerID int) bool {
	var exists bool
	database.Database.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM posts p WHERE p.id = ? AND `+VisibleCondition+`)`,
		postID, viewerID, viewerID, viewerID, viewerID,
	).Scan(&exists)
	return exists
}

func GetPostsByUser(userID, viewerID int) ([]models.Post, error) {
	query := PostSelect + ReactionSelect + UserReactionSelect + FromPosts + " WHERE p.user_id = ? AND p.group_id IS NULL AND " + VisibleCondition + " ORDER BY p.id DESC"

	rows, err := database.Database.Query(query, viewerID, userID, viewerID, viewerID, viewerID, viewerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []models.Post{}

	for rows.Next() {
		var p models.Post

		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Image, &p.Privacy, &p.GroupID, &p.UserID, &p.FullName, &p.Avatar, &p.CreatedAt, &p.LikeCount, &p.DislikeCount, &p.UserReaction); err != nil {
			return nil, err
		}

		result = append(result, p)
	}

	return result, rows.Err()
}

func GetPostsByGroup(groupID, viewerID int) ([]models.Post, error) {
	query := PostSelect + ReactionSelect + UserReactionSelect + FromPosts + " WHERE p.group_id = ? ORDER BY p.id DESC"

	rows, err := database.Database.Query(query, viewerID, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []models.Post{}

	for rows.Next() {
		var p models.Post

		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Image, &p.Privacy, &p.GroupID, &p.UserID, &p.FullName, &p.Avatar, &p.CreatedAt, &p.LikeCount, &p.DislikeCount, &p.UserReaction); err != nil {
			return nil, err
		}

		result = append(result, p)
	}

	return result, rows.Err()
}
