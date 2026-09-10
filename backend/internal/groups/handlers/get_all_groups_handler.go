package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"zone/internal/auth"
	"zone/internal/groups/repositories"
	"zone/internal/utils"
)

func GetAllGroupsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HandleError(w, 405, http.StatusText(405))
		return
	}

	user_id, err := auth.GetUserID(r)
	if err != nil {
		utils.HandleError(w, 401, http.StatusText(401))
		return
	}

	groups, err := repositories.GetAllGroups(user_id)
	if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}

	utils.RespondJSON(w, http.StatusOK, groups)
}

// url : /groups/{group_id}

func GetGroupInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.HandleError(w, 405, http.StatusText(405))
		return
	}
	_, err := auth.GetUserID(r)
	if err != nil {
		utils.HandleError(w, 401, http.StatusText(401))
		return
	}

	parts := strings.Split(r.URL.Path, "/")

	if len(parts) < 3 {
		utils.HandleError(w, 400, http.StatusText(400))
		return
	}

	group_id, err := strconv.Atoi(parts[2])
	if err != nil {
		utils.HandleError(w, 400, http.StatusText(400))
		return
	}

	groupinfo , err  := repositories.GetGroupbyid(group_id) 
		if err != nil {
		utils.HandleError(w, 500, http.StatusText(500))
		return
	}

	utils.RespondJSON(w, http.StatusOK, groupinfo)
}

