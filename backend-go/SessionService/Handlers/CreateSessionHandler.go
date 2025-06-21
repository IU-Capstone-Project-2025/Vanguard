package Handlers

import (
	"encoding/json"
	"net/http"
)

func (h *SessionManagerHandler) CreateSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.logger.Info("CreateSessionHandler", "Request Method", r.Method)
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}
	UserId := r.URL.Query().Get("userId")
	session, err := h.Manager.NewSession()
	AdminToken := h.Manager.GenerateUserToken(session.Code, UserId, "Admin")
	if err != nil {
		h.logger.Error("error With CreateSession", "CreateSessionHandler", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(AdminToken)
	if err != nil {
		h.logger.Error("CreateSessionHandler", "CreateSessionHandler", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
