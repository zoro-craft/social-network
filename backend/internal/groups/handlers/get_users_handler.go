package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"zone/internal/auth"
	"zone/internal/groups/repositories"
	"zone/internal/utils"
)

// this handler is used to get all users that are not members of a specific group.
func GetAllUsersNotmemberHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user_id, err := auth.GetUserID(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, http.StatusText(401))
		return
	}
	// get the grouidp id from the url
	parts := strings.Split(r.URL.Path, "/")

	if len(parts) < 4 {
		utils.HandleError(w, http.StatusBadRequest, http.StatusText(400))
		return
	}

	group_id, err := strconv.Atoi(parts[2])
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Invalid group ID")
		return
	}

	exists := repositories.IsMember(user_id, group_id)
	if !exists {
		utils.HandleError(w, 403, http.StatusText(403))
		return
	}
	// get all users that are not members of the group
	users, err := repositories.GetAllUsersNotmember(user_id, group_id)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Could not load users")
		return
	}

	utils.RespondJSON(w, http.StatusOK, users)
}

// GetGroupMembersHandler handles the request to retrieve all members of a specific group.
func GetGroupMembersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	user_id, err := auth.GetUserID(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, http.StatusText(401))
		return
	}

	parts := strings.Split(r.URL.Path, "/")

	if len(parts) < 4 {
		utils.HandleError(w, http.StatusBadRequest, http.StatusText(400))
		return
	}

	group_id, err := strconv.Atoi(parts[2])
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Invalid group ID")
		return
	}

	exists := repositories.IsMember(user_id, group_id)
	if !exists {
		utils.HandleError(w, 403, http.StatusText(403))
		return
	}
	// get all members of the group
	members, err := repositories.GetAllMembers(group_id)
	if err != nil {
		fmt.Println(err)
		utils.HandleError(w, http.StatusInternalServerError, "Could not load members")
		return
	}

	utils.RespondJSON(w, http.StatusOK, members)
}
