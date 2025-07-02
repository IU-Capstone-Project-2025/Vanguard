package Handlers

import (
	"encoding/json"
	"net/http"
	"xxx/SessionService/models"
)

func DeleteUserHandler(registry *ConnectionRegistry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO:
		//   - If token is missing or invalid, respond with appropriate HTTP status (e.g., 400/401) and return early.
		//   - On extractTokenData error, write an HTTP error or close the connection instead of only logging.
		//   - On upgrader.Upgrade error, write log and return so no further processing.
		//   - Ensure that after Upgrade, if token parsing failed, the connection is closed.

		// Extracts the "token" from URL query. If missing, it should reject the request
		tokenId := r.URL.Query().Get("code")
		tokenName := r.URL.Query().Get("userId")
		// Parses and validates the token via extractTokenData. If invalid or expired, it should reject the request
		flag := registry.UnregisterConnection(tokenId, tokenName)
		if !flag {
			registry.logger.Error("DeleteUserFromSessionHandler error to unregister connection",
				"userId", tokenId,
				"UserName", tokenName,
				"flag", flag,
			)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{Message: "err to unregister connection"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}
