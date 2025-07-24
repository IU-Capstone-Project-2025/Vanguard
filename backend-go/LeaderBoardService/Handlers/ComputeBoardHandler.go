package Handlers

import (
	"encoding/json"
	"net/http"
	"xxx/LeaderBoardService/models"
	"xxx/shared"
)

// ComputeBoardHandler godoc
// @Summary      Compute leaderboard and get popular answers
// @Description  Accepts session answers, computes user scores and returns leaderboard data
// @Tags         leaderboard
// @Accept       json
// @Produce      json
// @Param        data  body      shared.SessionAnswers  true  "Session Answers Payload"
// @Success      200   {object}  shared.BoardResponse
// @Failure      400   {object}  models.ErrorResponse   "Bad Request or computation error"
// @Failure      405   {object}  models.ErrorResponse   "Method Not Allowed"
// @Failure      500   {object}  models.ErrorResponse   "Internal Server Error"
// @Router /get-results [post]
func (m *HandlerManager) ComputeBoardHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		w.Header().Set("Content-Type", "application/json")
		m.log.Error("Only POST method is allowed ", "Request Method", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req shared.SessionAnswers
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		m.log.Error("ComputeBoardHandler err to decode req",
			"Decode err", err,
			"Request Body", r.Body)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Message: "Bad Request"})
		return
	}
	userScore, err := m.Service.ComputeLeaderBoard(req)
	if err != nil {
		m.log.Error("ComputeBoardHandler err to compute userScore", "err", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
	}
	ans, err := m.Service.PopularAns(req)
	if err != nil {
		m.log.Error("ComputeBoardHandler err to popular ans", "err", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	resp := shared.BoardResponse{
		SessionCode: req.SessionCode,
		Table:       userScore,
		Popular:     ans,
	}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		m.log.Error("ComputeBoardHandler err to write response",
			"response", resp,
			"err", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Message: "StatusInternalServerError"})
		return
	}
}
