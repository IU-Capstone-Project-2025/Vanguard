package Handlers

import (
	"encoding/json"
	"net/http"
	models2 "xxx/SessionService/models"
	"xxx/shared"
)

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
		json.NewEncoder(w).Encode(models2.ErrorResponse{Message: "Bad Request"})
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
		Table:   userScore,
		Popular: ans,
	}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		m.log.Error("ComputeBoardHandler err to write response",
			"response", resp,
			"err", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models2.ErrorResponse{Message: "StatusInternalServerError"})
		return
	}
}
