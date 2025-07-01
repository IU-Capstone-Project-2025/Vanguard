package Handlers

import (
	"encoding/json"
	"net/http"
	"xxx/SessionService/models"
)

// StartSessionHandler starts an existing session by its ID.
//
// @Summary Start a session
// @Description Starts a session using the provided session ID.
// @Tags sessions
// @Accept  json
// @Produce  json
// @Param   id   query    string  true  "Session ID"
// @Success 200 "Session started successfully"
// @Failure 405 {object} models.ErrorResponse "Method not allowed"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /start [post]
func (h *SessionManagerHandler) StartSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.ErrorResponse{Message: "Method not allowed"})
		return
	}

	req := r.URL.Query().Get("id")
	code := r.URL.Query().Get("code")
	err := h.Manager.SessionStart(req, code)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
}
