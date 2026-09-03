package services

import "zone/internal/groups/repositories"

func JoinRequestValidator(group_id int , user_id int) bool {
	if (!repositories.GroupExists(group_id)) {
		return false
	}
	if (repositories.Alreadymember(group_id , user_id)) {
		return  false
	}
	return true
}