package handlers

import (
	"net/http"

	"zone/internal/auth"
	"zone/internal/utils"
)

// ConversationsHandler serves GET /conversations: the whole sidebar in one response.
func (h *Handler) ConversationsHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserID(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	users, groups, err := h.svc.GetConversations(userID)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "failed to load conversations")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]any{
		"users":  users,
		"groups": groups,
	})
}