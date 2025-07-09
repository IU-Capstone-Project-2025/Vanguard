package Handlers

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"xxx/SessionService/models"
	"xxx/shared"
)

// CreateSessionHandlerMock creates a new session and generates an admin token for the user.
//
// @Summary Create a new session and get an admin token. Mock endpoint, no req to another service
// @Description Creates a new session and returns an admin token for the specified user by userId.
// @Tags sessions
// @Accept  json
// @Produce  json
// @Param   request  body  models.CreateSessionReq  true  " Create Session req"
// @Success 200 {object} models.SessionCreateResponse "Admin token in JSON format"
// @Failure 405 {object} models.ErrorResponse "Method not allowed, only GET is allowed"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /sessionsMock [post]
func (h *SessionManagerHandler) CreateSessionHandlerMock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.logger.Info("CreateSessionHandler request method not allowed ", "Request Method", r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.ErrorResponse{Message: "Only GET method is allowed"})
		return
	}
	var req models.CreateSessionReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("CreateSessionHandler err to decode req",
			"Decode err", err,
			"Request Body", r.Body)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Message: "Bad Request"})
		return
	}
	h.logger.Debug("CreateSessionHandler get req", "req", req)
	session, err := h.Manager.NewSession()
	AdminToken := h.Manager.GenerateUserToken(session.Code, req.UserName, shared.RoleAdmin)
	if err != nil {
		h.logger.Error("CreateSessionHandler Error with create Session", "err", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Message: "StatusInternalServerError"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = h.Manager.SessionStartMock(req.QuizId, AdminToken.SessionId)
	if err != nil {
		h.logger.Error("CreateSessionHandler error With SessionStart",
			"QuizId", req.QuizId,
			"SessionId", AdminToken.SessionId,
			"err", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Message: "StatusInternalServerError"})
		return
	}
	s := jwt.NewWithClaims(jwt.SigningMethodHS256, AdminToken)
	token, err := s.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		h.logger.Error("CreateSessionHandler error to generate jwt token for user",
			"AdminToken", AdminToken,
			"err", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Message: "StatusInternalServerError"})
		return
	}
	response := models.SessionCreateResponse{
		Jwt:              token,
		ServerWsEndpoint: shared.GetWsEndpoint(),
		SessionId:        AdminToken.SessionId,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		h.logger.Error("CreateSessionHandler err to write response",
			"response", response,
			"err", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Message: "StatusInternalServerError"})
		return
	}
	h.logger.Debug("CreateSessionHandler success encode response", "response", response)
}
