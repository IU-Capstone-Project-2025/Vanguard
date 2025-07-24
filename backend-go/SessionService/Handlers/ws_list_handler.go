package Handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"log/slog"
	"net/http"
	"sync"
	"time"
	"xxx/real_time/config"
	"xxx/shared"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  5 * 1024,
	WriteBufferSize: 5 * 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// ConnectionContext stores necessary data of user after successful WebSocket connection.
type ConnectionContext struct {
	Conn      *SafeConn // the connection tunnel with user
	UserId    string    // unique ID of the connected user
	UserName  string
	SessionId string          // session ID of the session user joined in
	Role      shared.UserRole // the role of the user within the session
}

// NewWebSocketHandler returns an http.HandlerFunc that uses the given registry.
// It handles WebSocket upgrade requests for real-time connections.
// It expects a query parameter "token" containing a valid JWT (ephemeral token).
//
// Extracts the "token" from URL query;
// Parses and validates the token;
// Upgrades the HTTP request to a WebSocket connection;
// Starts a goroutine running handleRead to process incoming messages
func NewWebSocketHandler(registry *ConnectionRegistry) http.HandlerFunc {
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
			registry.logger.Error("WsHandler error to extract token", "err", err)
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		// Upgrades the HTTP request to a WebSocket connection.
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			registry.logger.Error("WsHandler error to upgrade websocket", "err", err)
			return
		}

		// Ensure session exists in registry
		registry.RegisterSession(token.SessionId)

		// Register this connection
		ctx := &ConnectionContext{
			Conn:      &SafeConn{Conn: conn},
			UserId:    token.UserId,
			UserName:  token.UserName,
			SessionId: token.SessionId,
			Role:      token.UserType,
		}

		// Send a welcome message
		//welcome := fmt.Sprintf(`{"type":"welcome","sessionId":"%s","userId":"%s"}`, ctx.SessionId, ctx.UserId)
		//err = conn.WriteMessage(websocket.TextMessage, []byte(welcome))
		//if err != nil {
		//	registry.logger.Error("WsHandler error to send welcome", "err", err)
		//}
		//registry.logger.Info("WsHandler welcome to connection", "welcome", welcome)
		// Start reading messages for this connection in a separate goroutine.
		go handleRead(ctx, registry)
	}
}

type SafeConn struct {
	Conn  *websocket.Conn
	Mutex sync.Mutex
}

type ConnectionRegistry struct {
	mu          sync.RWMutex
	connections map[string]map[string]*SafeConn
	rooms       map[string]map[string]string
	logger      *slog.Logger
}

// NewConnectionRegistry initializes the ConnectionRegistry
func NewConnectionRegistry(log *slog.Logger) *ConnectionRegistry {
	return &ConnectionRegistry{
		connections: make(map[string]map[string]*SafeConn),
		rooms:       make(map[string]map[string]string),
		logger:      log,
	}
}

// RegisterSession creates a new session entry; idempotent
func (r *ConnectionRegistry) RegisterSession(sessionID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.connections[sessionID]; !exists {
		r.connections[sessionID] = make(map[string]*SafeConn)
		r.rooms[sessionID] = make(map[string]string)
	}
}

// UnregisterSession removes session entirely (e.g., on session end)
func (r *ConnectionRegistry) UnregisterSession(sessionID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.connections, sessionID)
	delete(r.rooms, sessionID)
}

// RegisterConnection adds new joined user connection, mapping to a corresponding session
func (r *ConnectionRegistry) RegisterConnection(sessionID, userID, userName string, conn *websocket.Conn) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	sessions, exists := r.connections[sessionID]
	if !exists {
		return fmt.Errorf("session %s not found", sessionID)
	}
	sessions[userID] = &SafeConn{Conn: conn}
	rooms, exists := r.rooms[sessionID]
	if !exists {
		return fmt.Errorf("session %s not found", sessionID)
	}
	rooms[userID] = userName
	return nil
}

// UnregisterConnection removes joined user connection, (e.g., on user disconnect)
func (r *ConnectionRegistry) UnregisterConnection(sessionID, userID string) bool {
	r.mu.Lock()
	e1 := false
	e2 := false
	if sessions, exists1 := r.connections[sessionID]; exists1 {
		e1 = true
		if sessions[userID] != nil && sessions[userID].Conn != nil {
			sessions[userID].Mutex.Lock()
			_ = sessions[userID].Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "user left"))
			sessions[userID].Mutex.Unlock()
			sessions[userID].Conn.Close()
		}
		delete(sessions, userID)
	}
	if rooms, exists2 := r.rooms[sessionID]; exists2 {
		e2 = true
		delete(rooms, userID)
	}
	if e1 == true && e2 == true {
		r.logger.Info("UnregisterConnection", "session", sessionID, "user", userID)
		r.mu.Unlock()
		handleDelete(sessionID, userID, r)
		return true
	}
	r.mu.Unlock()
	return false
}

// GetConnections gets a snapshot copy of connections to avoid holding lock during WriteMessage
func (r *ConnectionRegistry) GetConnections(sessionID string) []*SafeConn {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var conns []*SafeConn
	if sessions, exists := r.connections[sessionID]; exists {
		for _, c := range sessions {
			conns = append(conns, c)
		}
	}
	return conns
}

func (r *ConnectionRegistry) GetConnectionById(sessionID, UserId string) *SafeConn {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.connections[sessionID][UserId]
}

func (r *ConnectionRegistry) GetRooms(sessionID string) map[string]string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rooms := make(map[string]string)
	if sessions, exists := r.rooms[sessionID]; exists {
		for key, c := range sessions {
			rooms[key] = c
		}
	}
	return rooms
}

func extractTokenData(tokenString string) (*shared.UserToken, error) {
	token, err := jwt.ParseWithClaims(tokenString, &shared.UserToken{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.LoadConfig().JWT.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	switch {
	case token.Valid:
		fmt.Println("OK token")
	case errors.Is(err, jwt.ErrTokenMalformed):
		fmt.Errorf("malformed token: %w", err)
	default:
		fmt.Errorf("couldn't handle this token: %w", err)
	}

	claims, ok := token.Claims.(*shared.UserToken)
	fmt.Println("Decoded token: ", *claims)
	if !ok {
		return nil, fmt.Errorf("error decoding jwt")
	}

	return claims, nil
}

func handleRead(ctx *ConnectionContext, reg *ConnectionRegistry) {
	err := reg.RegisterConnection(ctx.SessionId, ctx.UserId, ctx.UserName, ctx.Conn.Conn)
	if err != nil {
		reg.logger.Error("WsHandler handleRead error to register connection", "UserId", ctx.UserId,
			"userName", ctx.UserName,
			"error", err)
		return
	}
	UserConn := reg.GetConnectionById(ctx.SessionId, ctx.UserId)
	UserConn.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	UserConn.Conn.SetPongHandler(func(string) error {
		UserConn.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			UserConn.Mutex.Lock()
			err := UserConn.Conn.WriteMessage(websocket.PingMessage, nil)
			UserConn.Mutex.Unlock()
			if err != nil {
				reg.logger.Info("Ping failed, closing connection", "userId", ctx.UserId)
				return
			}
			<-ticker.C
		}
	}()
	reg.logger.Info("ws connected to user", "userId", ctx.UserId, "userName", ctx.UserName)

	m := reg.GetRooms(ctx.SessionId)
	jsonData, err := json.Marshal(m)
	if err != nil {
		reg.logger.Error("WsHandler handleRead error to marshal json", err)
		return
	}
	for _, conn := range reg.GetConnections(ctx.SessionId) {
		conn.Mutex.Lock()
		err := conn.Conn.WriteMessage(websocket.TextMessage, jsonData)
		conn.Mutex.Unlock()
		if err != nil {
			reg.logger.Error("WsHandler handleRead error to write json", "err", err)
			continue
		}
	}
	for {
		_, _, err := UserConn.Conn.ReadMessage()
		if err != nil {
			reg.logger.Info("Connection closed", "userId", ctx.UserId, "err", err)
			break
		}
	}

	reg.UnregisterConnection(ctx.SessionId, ctx.UserId)
}

func handleDelete(sessionID, userID string, reg *ConnectionRegistry) {
	reg.logger.Info("ws admin delete user, or user close conn", "userId", userID, "SessionId", sessionID)
	fmt.Println("about to get rooms")
	rooms := reg.GetRooms(sessionID)
	fmt.Println("rooms gotten:", rooms)
	m := reg.GetRooms(sessionID)
	jsonData, err := json.Marshal(m)
	if err != nil {
		reg.logger.Error("WsHandler handleRead error to marshal json", err)
		return
	}
	for _, conn := range reg.GetConnections(sessionID) {
		conn.Mutex.Lock()
		err := conn.Conn.WriteMessage(websocket.TextMessage, jsonData)
		conn.Mutex.Unlock()
		if err != nil {
			reg.logger.Error("WsHandler handleRead error to write json", "err", err,
				"userId", userID,
				"SessionId", sessionID,
				"data", jsonData,
			)
			continue
		}
		reg.logger.Info("ws send message to user", "userId", userID)
	}
}
