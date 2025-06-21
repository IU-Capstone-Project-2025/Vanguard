package Handlers

import (
	"encoding/json"
	"net/http"
	"xxx/SessionService/models"
)

// ValidateCodeHandler validates a session code and returns a user token if valid.
//
// @Summary Validate session code
// @Description Validates a session code and returns a user token for the specified user if the code is valid.
// @Tags sessions
// @Accept  json
// @Produce  json
// @Param   code    query    string  true  "Session code"
// @Param   userId  query    string  true  "User ID"
// @Success 200 {object} models.UserToken "User token in JSON format"
// @Failure 400 {object} models.ErrorResponse "Invalid code"
// @Failure 405 {object} models.ErrorResponse "Only GET method is allowed"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /validate [get]
func (h *SessionManagerHandler) ValidateCodeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.ErrorResponse{Message: "Only GET method is allowed"})
		return
	}

	code := r.URL.Query().Get("code")
	userId := r.URL.Query().Get("userId")

	if h.Manager.ValidateCode(code) {
		userToken := h.Manager.GenerateUserToken(code, userId, "User")
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
