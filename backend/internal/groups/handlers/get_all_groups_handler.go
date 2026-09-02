package handlers

import (
	"net/http"

	"zone/internal/auth"
	"zone/internal/groups/repositories"
	"zone/internal/utils"
)

func GetAllGroupsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HandleError(w, 405, http.StatusText(405))
		return
	}

	userID, err := auth.GetUserIDFromSession(r)
	if err != nil {
		utils.HandleError(w, 401, http.StatusText(401))
		return
	}

	groups, err := repositories.GetAllGroupsNotMember(userID)
	if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}
	utils.RespondJSON(w, http.StatusOK, groups)
}
