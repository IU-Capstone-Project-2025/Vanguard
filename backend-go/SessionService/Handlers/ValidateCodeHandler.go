package Handlers

import (
	"encoding/json"
	"net/http"
	"xxx/SessionService/models"
	"xxx/shared"
)

// ValidateCodeHandler validates a session code and returns a user token if valid.
//
// @Summary Validate session code
// @Description Validates a session code and returns a user token for the specified user if the code is valid.
// @Tags sessions
// @Accept  json
// @Produce  json
// @Param   request  body  models.ValidateCodeReq  true  " Create Session req"
// @Success 200 {object} models.UserToken "User token in JSON format"
// @Failure 400 {object} models.ErrorResponse "Invalid code"
// @Failure 405 {object} models.ErrorResponse "Only GET method is allowed"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /join [post]
func (h *SessionManagerHandler) ValidateCodeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.ErrorResponse{Message: "Only GET method is allowed"})
		return
	}

	var req models.ValidateCodeReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Message: "Bad Request"})
		return
	}

	if h.Manager.ValidateCode(req.Code) {
		userToken := h.Manager.GenerateUserToken(req.Code, req.UserId, shared.RoleParticipant)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(userToken); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.ErrorResponse{Message: err.Error()})
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Message: "Code is incorrect"})
	}
}
