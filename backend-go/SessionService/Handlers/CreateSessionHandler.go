package Handlers

import (
	"encoding/json"
	"net/http"
	"xxx/SessionService/models"
)

// CreateSessionHandler creates a new session and generates an admin token for the user.
//
// @Summary Create a new session and get an admin token
// @Description Creates a new session and returns an admin token for the specified user by userId.
// @Tags sessions
// @Accept  json
// @Produce  json
// @Param   userId   query    string  true  "User ID"
// @Success 200 {object} models.UserToken "Admin token in JSON format"
// @Failure 405 {object} models.ErrorResponse "Method not allowed, only GET is allowed"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /create [get]
func (h *SessionManagerHandler) CreateSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.logger.Info("CreateSessionHandler", "Request Method", r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.ErrorResponse{Message: "Only GET method is allowed"})
	}
	UserId := r.URL.Query().Get("userId")
	session, err := h.Manager.NewSession()
	AdminToken := h.Manager.GenerateUserToken(session.Code, UserId, "Admin")
	if err != nil {
		h.logger.Error("error With CreateSession", "CreateSessionHandler", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Message: "StatusInternalServerError"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(AdminToken)
	if err != nil {
		h.logger.Error("CreateSessionHandler", "CreateSessionHandler", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Message: "StatusInternalServerError"})
		return
	}
}
