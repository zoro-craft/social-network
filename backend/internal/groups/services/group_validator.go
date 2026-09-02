package services

import (
	"strings"

	"zone/internal/groups/models"
)

// ValidateGroup checks the fields a client is allowed to provide when creating a group.
func ValidateGroup(groupe models.Groupe) bool {
	if (len(strings.TrimSpace(groupe.Titel)) <= 0 && len(strings.TrimSpace(groupe.Titel)) >= 10) {
		return  false
	}

	if (len(strings.TrimSpace(groupe.Description))<= 0 && len(strings.TrimSpace(groupe.Description)) >=200)   {
		return false
	}
	return true
}
