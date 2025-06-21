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
// @Param   code   query    string  true  "Session code"
// @Success 201 "Successfully moved to the next question"
// @Failure 405 {object} models.ErrorResponse "Method not allowed"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /session/{id}/nextQuestion [post]
func (h *SessionManagerHandler) NextQuestionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.ErrorResponse{Message: "Method not allowed"})
		return
	}

	vars := mux.Vars(r)
	code := vars["code"]
	err := h.Manager.NextQuestion(code)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Message: err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
}
