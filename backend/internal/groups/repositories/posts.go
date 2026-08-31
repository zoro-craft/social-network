package repositories

import (
	"database/sql"

	"zone/internal/database"
)

type GroupPost struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Image     string `json:"image,omitempty"`
	Nickname  string `json:"nickname,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}

type GroupComment struct {
	ID        int    `json:"id"`
	PostID    int    `json:"post_id"`
	UserID    int    `json:"user_id"`
	Content   string `json:"content"`
	Image     string `json:"image,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	Avatar    string `json:"avatar,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}

func CreateGroupPost(groupID, userID int, title, content, image string) (int64, error) {
	result, err := database.Database.Exec(
		"INSERT INTO posts (user_id, group_id, title, content, image, privacy) VALUES (?, ?, ?, ?, ?, 'private')",
		groupID, userID, title, content, image,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func GetGroupPost(postID int) (GroupPost, error) {
	p := GroupPost{}
	err := database.Database.QueryRow(`
		SELECT p.id, p.user_id, p.title, p.content, COALESCE(p.image, ''),
		       COALESCE(u.nickname, ''), p.created_at
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.id = ?`, postID,
	).Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.Image, &p.Nickname, &p.CreatedAt)
	return p, err
}

func ListGroupPosts(groupID int) ([]GroupPost, error) {
	rows, err := database.Database.Query(`
		SELECT p.id, p.user_id, p.title, p.content, COALESCE(p.image, ''),
		       COALESCE(u.nickname, ''), p.created_at
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.group_id = ?
		ORDER BY p.id DESC`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []GroupPost{}
	for rows.Next() {
		var p GroupPost
		if err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.Image, &p.Nickname, &p.CreatedAt); err == nil {
			posts = append(posts, p)
		}
	}
	return posts, rows.Err()
}

// GroupIDForPost returns the group_id a post belongs to, or 0 if it is not a group post.
func GroupIDForPost(postID int) (int, error) {
	var groupID int
	err := database.Database.QueryRow(
		"SELECT group_id FROM posts WHERE id = ?",
		postID,
	).Scan(&groupID)
	if err == sql.ErrNoRows {
		return 0, sql.ErrNoRows
	}
	return groupID, err
}

func CreateGroupComment(postID, userID, groupID int, content, image string) (int64, error) {
	result, err := database.Database.Exec(
		"INSERT INTO comments (post_id, user_id, group_id, content, image) VALUES (?, ?, ?, ?, ?)",
		postID, userID, groupID, content, image,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func ListGroupComments(postID int) ([]GroupComment, error) {
	rows, err := database.Database.Query(`
		SELECT c.id, c.post_id, c.user_id, c.content, COALESCE(c.image, ''),
		       COALESCE(u.first_name, ''), COALESCE(u.avatar, ''), c.created_at
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.post_id = ? AND c.group_id IS NOT NULL
		ORDER BY c.id ASC`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []GroupComment{}
	for rows.Next() {
		var c GroupComment
		if err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.Image, &c.FirstName, &c.Avatar, &c.CreatedAt); err == nil {
			comments = append(comments, c)
		}
	}
	return comments, rows.Err()
}
