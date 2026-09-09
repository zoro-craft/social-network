package handlers

import (
	"fmt"
	"net/http"

	"zone/internal/auth"
	"zone/internal/groups/repositories"
	"zone/internal/utils"
)

func GetAllGroupsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Println("1")
		utils.HandleError(w, 405, http.StatusText(405))
		return
	}

	userID, err := auth.GetUserID(r)
	if err != nil {
		fmt.Println("2")
		utils.HandleError(w, 401, http.StatusText(401))
		return
	}

	groups, err := repositories.GetAllGroupsNotMember(userID)
	if err != nil {
		fmt.Println("3")
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}
	utils.RespondJSON(w, http.StatusOK, groups)
}
