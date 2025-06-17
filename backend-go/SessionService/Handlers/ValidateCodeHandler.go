package Handlers

import (
	"encoding/json"
	"net/http"
)

func (h *SessionManagerHandler) ValidateCodeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}
	req := r.URL.Query().Get("code")
	flag := h.Manager.ValidateCode(req)
	if flag {
		userToken := h.Manager.GenerateToken(req)
		err := json.NewEncoder(w).Encode(userToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Code is incorrect", http.StatusBadRequest)
	}
}
