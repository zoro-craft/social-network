package repositories

import (
	"database/sql"
	"fmt"

	"zone/internal/database"
	"zone/internal/groups/models"
)

func SaveInvitation(tx *sql.Tx, invit models.GroupInvitation) (int, error) {
	var Id int
	// Check if invitation already exists
	err := tx.QueryRow(
		`
		SELECT id
		FROM group_invitations
		WHERE group_id = ?
		AND invited_user_id = ?
	`,
		invit.GroupID,
		invit.InvitedUserID,
	).Scan(&Id)

	if err == nil {
		return 0, fmt.Errorf("already invited")
	}

	result, err := tx.Exec(
		`
		INSERT INTO group_invitations
			(group_id, inviter_id, invited_user_id)
		VALUES (?, ?, ?)
	`,
		invit.GroupID,
		invit.InviterID,
		invit.InvitedUserID,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println("INSERT invitation error:", err)
		return 0, err
	}

	return int(id), nil
}

func InvitedUserIDExists(InvitedUserID int) bool {
	var exists bool
	err := database.Database.QueryRow(
		`SELECT EXISTS(
			SELECT 1 FROM users WHERE id = ?
		)`,
		InvitedUserID,
	).Scan(&exists)
	if err != nil {
		return false
	}

	return exists
}

func GetInvitation(invitationID, userID int) (models.GroupInvitation, error) {
	var invitation models.GroupInvitation

	err := database.Database.QueryRow(`
		SELECT id, group_id, inviter_id, invited_user_id, status, created_at
		FROM group_invitations
		WHERE id = ?
		AND invited_user_id = ?
	`, invitationID, userID).Scan(
		&invitation.ID,
		&invitation.GroupID,
		&invitation.InviterID,
		&invitation.InvitedUserID,
		&invitation.Status,
		&invitation.CreatedAt,
	)
	if err != nil {
		return invitation, err
	}

	return invitation, nil
}

// this function is used the get all users that are not members of a specific group .
func GetAllUsersNotmember(userid int, groupid int) ([]models.User, error) {
	rows, err := database.Database.Query(`
		SELECT id, full_name, COALESCE(avatar, '')
		FROM users
		WHERE id != ?
		AND id NOT IN (
			SELECT user_id
			FROM group_members
			WHERE group_id = ?
		)
		ORDER BY id
	`, userid, groupid)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var u models.User

		if err := rows.Scan(
			&u.ID,
			&u.Full_name,
			&u.Avatar,
		); err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, rows.Err()
}

// GetAllMembers retrieves all members of a specific group.
func GetAllMembers(groupID int) ([]models.User, error) {
	rows, err := database.Database.Query(`
		SELECT u.id, u.full_name, COALESCE(u.avatar, '')
		FROM users u
		JOIN group_members gm ON u.id = gm.user_id
		WHERE gm.group_id = ?
	`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var u models.User

		err := rows.Scan(
			&u.ID,
			&u.Full_name,
			&u.Avatar,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// GetInvitationIDByGroupAndUser returns the invitation id for a group and
// invited user, or sql.ErrNoRows when none exists.
func GetInvitationIDByGroupAndUser(groupID, invitedUserID int) (int, error) {
	var id int

	err := database.Database.QueryRow(`
		SELECT id
		FROM group_invitations
		WHERE group_id = ? AND invited_user_id = ?
	`, groupID, invitedUserID).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// delete the invitation in response
func DeleteInvitation(tx *sql.Tx, invitationID int) error {
    _, err := tx.Exec(`
        DELETE FROM group_invitations
        WHERE id = ?
    `, invitationID)

    return err
}
