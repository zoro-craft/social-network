package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"zone/internal/auth"
	"zone/internal/notifications/repositories"
	"zone/internal/notifications/services"
	"zone/internal/utils"
	"zone/internal/websocket"
)

// DefaultHandler is the shared handler wired to the real repository and realtime hub.
var DefaultHandler = NewHandler(services.NewService(repositories.NewNotificationRepository(), websocket.DefaultHub))

// Handler is the HTTP boundary for notifications: parse/auth, call the service, shape the response.
type Handler struct {
	svc *services.Service
}

// NewHandler builds a notifications HTTP handler.
func NewHandler(svc *services.Service) *Handler {
	return &Handler{svc: svc}
}

// GetNotifications serves GET /notifications: the recipient's own unread rows, newest first,
// 10 per page via an id keyset cursor, plus the exact total unread count.
func (h *Handler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserID(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	var cursorID int
	if q := r.URL.Query().Get("before_id"); q != "" {
		cursorID, _ = strconv.Atoi(q)
	}

	items, hasMore, err := h.svc.ListUnread(userID, cursorID)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "failed to load notifications")
		return
	}
	unreadCount, err := h.svc.CountUnread(userID)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "failed to load unread count")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]any{
		"items":        items,
		"unread_count": unreadCount,
		"has_more":     hasMore,
	})
}


func (h *Handler) MarkReadByID(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserID(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		utils.HandleError(w, http.StatusBadRequest, "invalid notification id")
		return
	}

	err = h.svc.MarkReadByID(userID, id)
	if errors.Is(err, sql.ErrNoRows) {
		utils.HandleError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "failed to mark notification read")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "read"})
}
