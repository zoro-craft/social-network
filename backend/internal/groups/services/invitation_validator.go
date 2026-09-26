package services

import (
	"zone/internal/groups/models"
	"zone/internal/groups/repositories"
)

func ValidCreateInvitation(invt models.GroupInvitation, userID int) bool {
	//check if the group exists
	if !repositories.GroupExists(invt.GroupID) {
		return false
	}
	// check if the invited user exists
	if !repositories.InvitedUserIDExists(invt.InvitedUserID) {
		return false
	}

	// check if the inviter is a member of the group
	if !repositories.IsMember(userID, invt.GroupID) {
		return false
	}
	// check if the invited user is already a member of the group
	if repositories.IsMember(invt.InvitedUserID, invt.GroupID) {
		return false
	}
	// check if the inviter is the same as the invited user
	if invt.InviterID != userID {
		return false
	}
	// check if the inviter is trying to invite themselves
	if invt.InvitedUserID == userID {
		return false
	}

	return true
}
