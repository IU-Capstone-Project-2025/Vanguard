package Handlers

import "net/http"

func (h *SessionManagerHandler) StartSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	req := r.URL.Query().Get("id")
	err := h.Manager.SessionStart(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	return
}
