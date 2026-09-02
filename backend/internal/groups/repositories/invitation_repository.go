package repositories

import (
	"zone/internal/database"
	"zone/internal/groups/models"
)

func SaveInvitation(invit models.GroupInvitation) (int, error) {
	result, err := database.Database.Exec(`
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


func CreateNotification(userID int, invitationID int , msg string) error {
	_, err := database.Database.Exec(`
		INSERT INTO notifications
			(user_id, type, invitation_id , message)
		VALUES (?, ?, ?, ?)
	`,
		userID,
		"group_invitation",
		invitationID,
		msg,
	)

	return err
}



func InvitationExists(invitationID, userID int) (bool) {
	var exists bool

	err := database.Database.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM group_invitations
			WHERE id = ?
			AND invited_user_id = ?
		)
	`, invitationID, userID).Scan(&exists)

	if err != nil {
		return false
	}

	return exists
}