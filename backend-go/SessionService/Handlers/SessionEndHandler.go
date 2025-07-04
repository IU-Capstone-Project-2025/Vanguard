package Handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"xxx/SessionService/models"
)

// SessionEndHandler used to end session.
//
// @Summary delete session, send message to rabbit
// @Description Delete session from redis, send message to rabbit that session deleted
// @Tags sessions
// @Accept  json
// @Produce  json
// @Param   id   path   string  true  "Session ID"
// @Success 200 "Successfully moved to the next question"
// @Failure 405 {object} models.ErrorResponse "Method not allowed"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /session/{id}/end [post]
func (h *SessionManagerHandler) SessionEndHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.logger.Info("SessionEndHandler request method not allowed ", "Request Method", r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.ErrorResponse{Message: "Method not allowed"})
		return
	}
	vars := mux.Vars(r)
	code := vars["id"]
	err := h.Manager.SessionEnd(code)
	if err != nil {
		h.logger.Info("SessionEndHandler error to send next Question message to rabbit",
			"code", code,
			"err", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Message: err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	h.logger.Info("SessionEndHandler success", "code", code)
}
