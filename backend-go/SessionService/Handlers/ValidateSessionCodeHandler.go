package Handlers

import (
	"encoding/json"
	"net/http"
	"xxx/SessionService/models"
)

// ValidateSessionCodeHandler validates a session code and returns a user token if valid.
//
// @Summary Validate session code
// @Description Validates a session code and returns a user token for the specified user if the code is valid.
// @Tags sessions
// @Accept  json
// @Produce  json
// @Param   request  body  models.ValidateCodeReq  true  " Create Session req"
// @Success 200 {object} models.SessionCreateResponse "User token in JSON format"
// @Failure 400 {object} models.ErrorResponse "Invalid code"
// @Failure 405 {object} models.ErrorResponse "Only GET method is allowed"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /join [post]
func (h *SessionManagerHandler) ValidateSessionCodeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.logger.Info("ValidateCodeHandler request method not allowed ", "Request Method", r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.ErrorResponse{Message: "Only GET method is allowed"})
		return
	}
	var req models.ValidateSessionCodeReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("ValidateCodeHandler Request Body Decode Error",
			"body", r.Body,
			"Error", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Message: "Bad Request"})
		return
	}
	h.logger.Debug("ValidateCodeHandler Request Body", "Body", req)
	if h.Manager.ValidateCode(req.Code) {
		h.logger.Info("ValidateSessionCode code exist", "code", req.Code)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	} else {
		h.logger.Error("ValidateCodeHandler err to validate code", "code", req.Code)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Message: "Code is incorrect"})
		return
	}
}
