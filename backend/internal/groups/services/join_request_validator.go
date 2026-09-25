package services

import (

	"zone/internal/groups/repositories"
)

func JoinRequestValidator(group_id int, user_id int) bool {
	if !repositories.GroupExists(group_id) {
		return false
	}
	if repositories.IsMember( user_id , group_id) {
		return false
	}

	return true
}
