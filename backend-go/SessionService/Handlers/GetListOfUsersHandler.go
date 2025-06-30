package Handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"xxx/SessionService/models"
)

// GetListOfUsers sends list of users of current session
//
// @Summary sends list of users
// @Description
// @Tags sessions
// @Accept  json
// @Produce  json
// @Param   id   path   string  true  "id"
// @Success 200 {object} models.GetPlayersResponseGetPlayersResponse ""
// @Failure 405 {object} models.ErrorResponse "Method not allowed, only GET is allowed"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /session/{id}/list [post]
func (h *SessionManagerHandler) GetListOfUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.logger.Info("GetListOfUsers request method not allowed ", "Request Method", r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.ErrorResponse{Message: "Only GET method is allowed"})
		return
	}
	vars := mux.Vars(r)
	code := vars["id"]
	users, err := h.Manager.GetListOfUsers(code)
	if err != nil {
		h.logger.Info("GetListOfUsers error to get list of students",
			"code", code,
			"err", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Message: err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response := models.GetPlayersResponse{
		SessionCode: code,
		Players:     users,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		h.logger.Error("GetListOfUsers err to write response",
			"response", response,
			"err", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Message: "StatusInternalServerError"})
		return
	}
}
