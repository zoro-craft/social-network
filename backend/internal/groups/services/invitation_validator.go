package services

import (
	"zone/internal/groups/models"
	"zone/internal/groups/repositories"
)

func ValidCreateInvitation(invt models.GroupInvitation, userID int) bool {
	if !repositories.GroupExists(invt.GroupID) {
		return false
	}

	if !repositories.InvitedUserIDExists(invt.InvitedUserID) {
		return false
	}
	if !repositories.IsMember(userID, invt.GroupID) {
		return false
	}
	if invt.InviterID != userID {
		return false
	}

	return true
}
