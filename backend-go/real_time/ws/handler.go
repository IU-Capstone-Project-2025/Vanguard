package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"xxx/shared"
)

type HandlerDeps struct {
	Tracker  *QuizTracker
	Registry *ConnectionRegistry
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  5 * 1024,
	WriteBufferSize: 5 * 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// ConnectionContext stores necessary data of user after successful WebSocket connection.
type ConnectionContext struct {
	Conn      *websocket.Conn // the connection tunnel with user
	UserId    string          // unique ID of the connected user
	SessionId string          // session ID of the session user joined in
	Role      shared.UserRole // the role of the user within the session
	mu        sync.Mutex
}

// NewWebSocketHandler returns a http.HandlerFunc that uses the given registry.
// It handles WebSocket upgrade requests for real-time connections.
// It expects a query parameter "token" containing a valid JWT (ephemeral token).
//
// Extracts the "token" from URL query;
// Parses and validates the token;
// Upgrades the HTTP request to a WebSocket connection;
// Starts a goroutine running handleRead to process incoming messages
func NewWebSocketHandler(deps HandlerDeps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO:
		//   - If token is missing or invalid, respond with appropriate HTTP status (e.g., 400/401) and return early.
		//   - On extractTokenData error, write an HTTP error or close the connection instead of only logging.
		//   - On upgrader.Upgrade error, write log and return so no further processing.
		//   - Ensure that after Upgrade, if token parsing failed, the connection is closed.

		// Extracts the "token" from URL query. If missing, it should reject the request
		tokenString := r.URL.Query().Get("token")
		if tokenString == "" {
			http.Error(w, "missing token", http.StatusBadRequest)
			return
		}

		// Parses and validates the token via extractTokenData. If invalid or expired, it should reject the request
		token, err := extractTokenData(tokenString)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		// Upgrades the HTTP request to a WebSocket connection.
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("ws upgrade error:", err)
			return
		}
		fmt.Println("Try to register connection in handler.go")

		// Ensure session exists in registry
		deps.Registry.RegisterSession(token.SessionId)

		// Register this connection
		ctx := &ConnectionContext{
			Conn:      conn,
			UserId:    token.UserId,
			SessionId: token.SessionId,
			Role:      token.UserType,
			mu:        sync.Mutex{},
		}
		fmt.Println("Try to register user in handler.go")
		if err := deps.Registry.RegisterConnection(ctx); err != nil {
			log.Printf("Failed to register connection: %v", err)
			conn.Close()
			return
		}

		// Send a welcome message
		welcome := fmt.Sprintf(`{"type":"welcome","sessionId":"%s","userId":"%s"}`, ctx.SessionId, ctx.UserId)
		conn.WriteMessage(websocket.TextMessage, []byte(welcome))

		// Start reading messages for this connection in a separate goroutine.
		go handleRead(ctx, deps)
	}
}
