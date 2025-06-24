package Handlers

import (
	"encoding/json"
	"net/http"
	"xxx/SessionService/models"
	"xxx/shared"
)

// CreateSessionHandler creates a new session and generates an admin token for the user.
//
// @Summary Create a new session and get an admin token
// @Description Creates a new session and returns an admin token for the specified user by userId.
// @Tags sessions
// @Accept  json
// @Produce  json
// @Param   request  body  models.CreateSessionReq  true  " Create Session req"
// @Success 200 {object} models.UserToken "Admin token in JSON format"
// @Failure 405 {object} models.ErrorResponse "Method not allowed, only GET is allowed"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /sessions [post]
func (h *SessionManagerHandler) CreateSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.logger.Info("CreateSessionHandler", "Request Method", r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.ErrorResponse{Message: "Only GET method is allowed"})
	}
	var req models.CreateSessionReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Message: "Bad Request"})
		return
	}
	session, err := h.Manager.NewSession()
	AdminToken := h.Manager.GenerateUserToken(session.Code, req.UserId, shared.RoleAdmin)
	if err != nil {
		h.logger.Error("error With CreateSession", "CreateSessionHandler", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Message: "StatusInternalServerError"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = h.Manager.SessionStart(req.UserId)
	if err != nil {
		h.logger.Error("error With CreateSession", "CreateSessionHandler", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Message: "StatusInternalServerError"})
	}
	err = json.NewEncoder(w).Encode(AdminToken)
	if err != nil {
		h.logger.Error("CreateSessionHandler", "CreateSessionHandler", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Message: "StatusInternalServerError"})
		return
	}
}
