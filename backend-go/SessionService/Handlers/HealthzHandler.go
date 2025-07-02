package Handlers

import (
	"encoding/json"
	"net/http"
	"xxx/SessionService/models"
)

func (h *SessionManagerHandler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.logger.Info("CreateSessionHandler request method not allowed ", "Request Method", r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.ErrorResponse{Message: "Only GET method is allowed"})
		return
	}
	err := h.Manager.CheckService()
	if err != nil {
		h.logger.Error("HealthHandler err",
			"err", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Message: "Internal Server Error"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
