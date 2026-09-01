package repositories

import (
	"zone/internal/database"
	"zone/internal/groups/models"
)

func GetInvitationsByUserID(userID int) ([]models.GroupInvitation, error) {
	rows, err := database.Database.Query(`
		SELECT id, group_id, inviter_id, invited_user_id, status, created_at
		FROM group_invitations
		WHERE invited_user_id = ?
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invitations []models.GroupInvitation

	for rows.Next() {
		var invitation models.GroupInvitation

		err := rows.Scan(
			&invitation.ID,
			&invitation.GroupID,
			&invitation.InviterID,
			&invitation.InvitedUserID,
			&invitation.Status,
			&invitation.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		invitations = append(invitations, invitation)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return invitations, nil
}
