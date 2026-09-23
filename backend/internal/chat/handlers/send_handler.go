package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"zone/internal/auth"
	"zone/internal/chat/repositories"
	"zone/internal/chat/services"
	"zone/internal/utils"
	"zone/internal/websocket"
)

// DefaultHandler is the shared chat handler wired to the real repository and realtime hub.
var DefaultHandler = NewHandler(services.NewService(repositories.NewChatRepository(), websocket.DefaultHub))

// Handler is the HTTP boundary for chat: parse/auth, call the service, shape the response.
type Handler struct {
	svc *services.Service
}

// NewHandler builds a chat HTTP handler.
func NewHandler(svc *services.Service) *Handler {
	return &Handler{svc: svc}
}

type sendDMRequest struct {
	RecipientID int    `json:"recipient_id"`
	Content     string `json:"content"`
}

// SendDMHandler serves POST /messages: send a direct message.
func (h *Handler) SendDMHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserID(r)
	if err != nil {
		utils.HandleError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	var req sendDMRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.svc.SendDM(userID, req.RecipientID, req.Content); err != nil {
		h.writeSendError(w, err)
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]any{})
}

type sendGroupRequest struct {
	Content string `json:"content"`
}

// SendGroupHandler serves POST /groups/{group_id}/messages: send a group message.
func (h *Handler) SendGroupHandler(w http.ResponseWriter, r *http.Request) {
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

	var req sendGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.HandleError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.svc.SendGroup(userID, groupID, req.Content); err != nil {
		h.writeSendError(w, err)
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]any{})
}

// writeSendError maps a send rejection to the agreed contract response.
func (h *Handler) writeSendError(w http.ResponseWriter, err error) {
	if errors.Is(err, services.ErrContentLong) || errors.Is(err, services.ErrNotAllowed) || errors.Is(err, services.ErrNotMember) {
		utils.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.HandleError(w, http.StatusInternalServerError, "failed to send message")
}