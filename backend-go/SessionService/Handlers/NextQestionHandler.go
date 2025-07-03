package Handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"xxx/SessionService/models"
)

// NextQuestionHandler advances to the next question for the given session code.
//
// @Summary Move to the next question
// @Description Advances to the next question in the session identified by the provided code.
// @Tags sessions
// @Accept  json
// @Produce  json
// @Param   id   path   string  true  "Session ID"
// @Success 200 "Successfully moved to the next question"
// @Failure 405 {object} models.ErrorResponse "Method not allowed"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /session/{id}/nextQuestion [post]
func (h *SessionManagerHandler) NextQuestionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.logger.Info("NextQuestionHandler request method not allowed ", "Request Method", r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.ErrorResponse{Message: "Method not allowed"})
		return
	}
	vars := mux.Vars(r)
	code := vars["id"]
	err := h.Manager.NextQuestion(code)
	if err != nil {
		h.logger.Info("NextQuestionHandler error to send next Question message to rabbit",
			"code", code,
			"err", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Message: err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	h.logger.Info("NextQuestionHandler success", "code", code)
}
