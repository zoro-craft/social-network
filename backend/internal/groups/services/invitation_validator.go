package services

import (
	"zone/internal/groups/models"
	"zone/internal/groups/repositories"
)

func ValidCreateInvitation(invt models.GroupInvitation, userID int) bool {
	if (!repositories.GroupExists(invt.GroupID)) {
		return  false
	}
	if invt.InviterID <= 0 {
		return false
	}

	if (!repositories.InvitedUserIDExists(invt.InvitedUserID)) {
		return false
	}

	if invt.InviterID != userID {
		return false
	}

	// Only pending invitations can be accepted/rejected
	if invt.Status != "pending" {
		return false
	}

	return true
}

