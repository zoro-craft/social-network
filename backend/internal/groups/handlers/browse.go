package handlers

import (
	"net/http"
	"strconv"

	"zone/internal/auth"
	"zone/internal/groups/repositories"
	"zone/internal/groups/services"
	"zone/internal/utils"
)

// BrowseGroupsHandler lists all available groups with membership info for the
// current user.
func BrowseGroupsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, _ := auth.GetUserIDFromSession(r)

	groups, err := repositories.ListGroups()
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to load groups")
		return
	}

	for i := range groups {
		groups[i].IsCreator = groups[i].Creator == userID
		groups[i].IsMember = services.IsGroupMember(groups[i].ID, userID) || groups[i].IsCreator
	}

	utils.RespondJSON(w, http.StatusOK, groups)
}

// GetGroupHandler returns a single group's details, only visible details. Group
// content (posts/events) are handled separately and require membership; the
// group info itself (title, description) is public for browsing.
func GetGroupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	groupID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || groupID <= 0 {
		utils.HandleError(w, http.StatusBadRequest, "Invalid group id")
		return
	}

	userID, _ := auth.GetUserIDFromSession(r)

	group, err := repositories.GetGroup(groupID)
	if err != nil {
		utils.HandleError(w, http.StatusNotFound, "Group not found")
		return
	}

	group.IsCreator = group.Creator == userID
	group.IsMember = services.IsGroupMember(groupID, userID) || group.IsCreator

	utils.RespondJSON(w, http.StatusOK, group)
}

// ListMembersHandler returns all members of a group. Only members can view it.
func ListMembersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, err := auth.GetUserIDFromSession(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, "You must be logged in")
		return
	}

	groupID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || groupID <= 0 {
		utils.HandleError(w, http.StatusBadRequest, "Invalid group id")
		return
	}

	if !services.GroupExists(groupID) {
		utils.HandleError(w, http.StatusNotFound, "Group not found")
		return
	}

	if !services.IsGroupMember(groupID, userID) && !services.IsGroupCreator(groupID, userID) {
		utils.HandleError(w, http.StatusForbidden, "You must be a member to view group members")
		return
	}

	members, err := repositories.ListMembers(groupID)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to load members")
		return
	}

	utils.RespondJSON(w, http.StatusOK, members)
}
