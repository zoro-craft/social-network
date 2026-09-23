package handlers

import (
	"net/http"
	"strconv"

	"zone/internal/auth"
	"zone/internal/utils"
)

// DMHistoryHandler serves GET /messages/{user_id}: direct conversation history.
func (h *Handler) DMHistoryHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserID(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	otherID, err := strconv.Atoi(r.PathValue("user_id"))
	if err != nil || otherID <= 0 {
		utils.HandleError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	h.serveHistory(w, r, userID, otherID, "private")
}

// GroupHistoryHandler serves GET /groups/{group_id}/messages: group history.
func (h *Handler) GroupHistoryHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserID(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	groupID, err := strconv.Atoi(r.PathValue("group_id"))
	if err != nil || groupID <= 0 {
		utils.HandleError(w, http.StatusBadRequest, "invalid group id")
		return
	}

	h.serveHistory(w, r, userID, groupID, "group")
}

// serveHistory builds the cursor-paged history response shared by both routes.
func (h *Handler) serveHistory(w http.ResponseWriter, r *http.Request, userID, targetID int, msgType string) {
	beforeID, _ := strconv.Atoi(r.URL.Query().Get("before"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	messages, hasMore, err := h.svc.ListHistory(userID, targetID, msgType, beforeID, limit)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]any{
		"messages": messages,
		"has_more": hasMore,
	})
}