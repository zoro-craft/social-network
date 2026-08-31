package repositories

import (
	"database/sql"
	"errors"

	"zone/internal/database"
	"zone/internal/groups/models"
)

func CreateGroup(title, description string, creatorID int) (int64, error) {
	result, err := database.Database.Exec(
		"INSERT INTO user_groups (creator_id, title, description) VALUES (?, ?, ?)",
		creatorID, title, description,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func GetGroup(groupID int) (models.Group, error) {
	g := models.Group{}
	err := database.Database.QueryRow(`
		SELECT g.id, g.creator_id, g.title, g.description, g.created_at,
		       COALESCE(u.first_name, '')
		FROM user_groups g
		JOIN users u ON g.creator_id = u.id
		WHERE g.id = ?`, groupID,
	).Scan(&g.ID, &g.Creator, &g.Title, &g.Description, &g.CreatedAt, &g.CreatorName)
	return g, err
}

func ListGroups() ([]models.Group, error) {
	rows, err := database.Database.Query(`
		SELECT g.id, g.creator_id, g.title, g.description, g.created_at,
		       COALESCE(u.first_name, ''),
		       (SELECT COUNT(*) FROM group_members m WHERE m.group_id = g.id)
		FROM user_groups g
		JOIN users u ON g.creator_id = u.id
		ORDER BY g.id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := []models.Group{}
	for rows.Next() {
		var g models.Group
		if err := rows.Scan(&g.ID, &g.Creator, &g.Title, &g.Description, &g.CreatedAt, &g.CreatorName, &g.MemberCount); err == nil {
			groups = append(groups, g)
		}
	}
	return groups, rows.Err()
}

func AddMember(groupID, userID int) error {
	_, err := database.Database.Exec(
		"INSERT OR IGNORE INTO group_members (group_id, user_id) VALUES (?, ?)",
		groupID, userID,
	)
	return err
}

func ListMembers(groupID int) ([]models.GroupMember, error) {
	rows, err := database.Database.Query(`
		SELECT m.id, m.group_id, m.user_id, m.joined_at,
		       COALESCE(u.first_name, ''), COALESCE(u.avatar, ''), COALESCE(u.nickname, '')
		FROM group_members m
		JOIN users u ON m.user_id = u.id
		WHERE m.group_id = ?
		ORDER BY m.joined_at ASC`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := []models.GroupMember{}
	for rows.Next() {
		var m models.GroupMember
		if err := rows.Scan(&m.ID, &m.GroupID, &m.UserID, &m.JoinedAt, &m.FirstName, &m.Avatar, &m.Nickname); err == nil {
			members = append(members, m)
		}
	}
	return members, rows.Err()
}

func HasInvitation(groupID, userID int) (bool, error) {
	var exists bool
	err := database.Database.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM group_invitations WHERE group_id = ? AND invited_user_id = ? AND status = 'pending')",
		groupID, userID,
	).Scan(&exists)
	return exists, err
}

func CreateInvitation(groupID, inviterID, invitedUserID int) error {
	_, err := database.Database.Exec(
		"INSERT OR IGNORE INTO group_invitations (group_id, inviter_id, invited_user_id, status) VALUES (?, ?, ?, 'pending')",
		groupID, inviterID, invitedUserID,
	)
	return err
}

func ListInvitationsForUser(userID int) ([]models.GroupInvitation, error) {
	rows, err := database.Database.Query(`
		SELECT i.id, i.group_id, i.inviter_id, i.invited_user_id, i.status, i.created_at,
		       g.title, COALESCE(g.description, ''), COALESCE(u.first_name, '')
		FROM group_invitations i
		JOIN user_groups g ON i.group_id = g.id
		JOIN users u ON g.creator_id = u.id
		WHERE i.invited_user_id = ? AND i.status = 'pending'
		ORDER BY i.created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	invitations := []models.GroupInvitation{}
	for rows.Next() {
		var i models.GroupInvitation
		if err := rows.Scan(&i.ID, &i.GroupID, &i.InviterID, &i.UserID, &i.Status, &i.CreatedAt, &i.GroupTitle, &i.GroupDesc, &i.CreatorName); err == nil {
			invitations = append(invitations, i)
		}
	}
	return invitations, rows.Err()
}

func AcceptInvitation(invitationID, userID int) error {
	return database.WithTx(func(tx *sql.Tx) error {
		var groupID int
		err := tx.QueryRow(
			"SELECT group_id FROM group_invitations WHERE id = ? AND invited_user_id = ? AND status = 'pending'",
			invitationID, userID,
		).Scan(&groupID)
		if err != nil {
			return errors.New("invitation not found")
		}

		if _, err := tx.Exec(
			"INSERT OR IGNORE INTO group_members (group_id, user_id) VALUES (?, ?)",
			groupID, userID,
		); err != nil {
			return err
		}

		_, err = tx.Exec(
			"DELETE FROM group_invitations WHERE id = ?",
			invitationID,
		)
		return err
	})
}

func RejectInvitation(invitationID, userID int) error {
	_, err := database.Database.Exec(
		"DELETE FROM group_invitations WHERE id = ? AND invited_user_id = ?",
		invitationID, userID,
	)
	return err
}

var ErrRequestNotFound = errors.New("request not found")

func HasPendingRequest(groupID, userID int) (bool, error) {
	var exists bool
	err := database.Database.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM group_requests WHERE group_id = ? AND user_id = ? AND status = 'pending')",
		groupID, userID,
	).Scan(&exists)
	return exists, err
}

func CreateRequest(groupID, userID int) error {
	_, err := database.Database.Exec(
		"INSERT OR IGNORE INTO group_requests (group_id, user_id, status) VALUES (?, ?, 'pending')",
		groupID, userID,
	)
	return err
}

func ListRequestsForGroup(groupID int) ([]models.GroupRequest, error) {
	rows, err := database.Database.Query(`
		SELECT r.id, r.group_id, r.user_id, r.status, r.created_at,
		       COALESCE(u.first_name, ''), COALESCE(u.nickname, ''), COALESCE(u.avatar, '')
		FROM group_requests r
		JOIN users u ON r.user_id = u.id
		WHERE r.group_id = ? AND r.status = 'pending'
		ORDER BY r.created_at ASC`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	requests := []models.GroupRequest{}
	for rows.Next() {
		var r models.GroupRequest
		if err := rows.Scan(&r.ID, &r.GroupID, &r.UserID, &r.Status, &r.CreatedAt, &r.FirstName, &r.Nickname, &r.Avatar); err == nil {
			requests = append(requests, r)
		}
	}
	return requests, rows.Err()
}

func AcceptRequest(requestID, userID int) error {
	return database.WithTx(func(tx *sql.Tx) error {
		var groupID int
		var requesterID int
		err := tx.QueryRow(
			"SELECT group_id, user_id FROM group_requests WHERE id = ? AND status = 'pending'",
			requestID,
		).Scan(&groupID, &requesterID)
		if err != nil {
			return ErrRequestNotFound
		}

		if _, err := tx.Exec(
			"INSERT OR IGNORE INTO group_members (group_id, user_id) VALUES (?, ?)",
			groupID, requesterID,
		); err != nil {
			return err
		}

		_, err = tx.Exec("DELETE FROM group_requests WHERE id = ?", requestID)
		return err
	})
}

func RefuseRequest(requestID int) error {
	result, err := database.Database.Exec(
		"DELETE FROM group_requests WHERE id = ? AND status = 'pending'",
		requestID,
	)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrRequestNotFound
	}
	return nil
}
