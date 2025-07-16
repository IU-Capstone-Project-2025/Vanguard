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
	Conn      *websocket.Conn // the connection tunnel with user
	UserId    string          // unique ID of the connected user
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
		if err := registry.RegisterConnection(token.SessionId, token.UserId, token.UserName, conn); err != nil {
			registry.logger.Error("WsHandler error to register connection",
				"SessionId", token.SessionId,
				"UserId", token.UserId,
				"UserName", token.UserName,
				"err", err)
			conn.Close()
			return
		}

		ctx := &ConnectionContext{
			Conn:      conn,
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

type ConnectionRegistry struct {
	mu          sync.RWMutex
	connections map[string]map[string]*websocket.Conn
	rooms       map[string]map[string]string
	logger      *slog.Logger
}

// NewConnectionRegistry initializes the ConnectionRegistry
func NewConnectionRegistry(log *slog.Logger) *ConnectionRegistry {
	return &ConnectionRegistry{
		connections: make(map[string]map[string]*websocket.Conn),
		rooms:       make(map[string]map[string]string),
		logger:      log,
	}
}

// RegisterSession creates a new session entry; idempotent
func (r *ConnectionRegistry) RegisterSession(sessionID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.connections[sessionID]; !exists {
		r.connections[sessionID] = make(map[string]*websocket.Conn)
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
	sessions[userID] = conn
	rooms, exists := r.rooms[sessionID]
	if !exists {
		return fmt.Errorf("session %s not found", sessionID)
	}
	rooms[userID] = userName
	go r.HandleCloseConnection(userID, userName, conn)
	return nil
}

// UnregisterConnection removes joined user connection, (e.g., on user disconnect)
func (r *ConnectionRegistry) UnregisterConnection(sessionID, userID string) bool {
	r.mu.Lock()
	e1 := false
	e2 := false
	if sessions, exists1 := r.connections[sessionID]; exists1 {
		e1 = true
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
func (r *ConnectionRegistry) GetConnections(sessionID string) []*websocket.Conn {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var conns []*websocket.Conn
	if sessions, exists := r.connections[sessionID]; exists {
		for _, c := range sessions {
			conns = append(conns, c)
		}
	}
	return conns
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

// новый участник = отправляем всем итоговую мапу
// админ удаляет чела или чел выходит = отправляем всем мапу
func handleRead(ctx *ConnectionContext, reg *ConnectionRegistry) {
	err := reg.RegisterConnection(ctx.SessionId, ctx.UserId, ctx.UserName, ctx.Conn)
	if err != nil {
		reg.logger.Error("WsHandler handleRead error to register connection")
		return
	}
	reg.logger.Info("ws connected to user", "userId", ctx.UserId, "userName", ctx.UserName)
	m := reg.GetRooms(ctx.SessionId)
	//con := reg.connections[ctx.SessionId][ctx.UserId]
	jsonData, err := json.Marshal(m)
	if err != nil {
		reg.logger.Error("WsHandler handleRead error to marshal json", err)
		return
	}
	//if len(m) > 0 {
	//	fmt.Println(m)
	//	err = con.WriteMessage(websocket.TextMessage, jsonData)
	//	if err != nil {
	//		reg.logger.Error("WsHandler handleRead error to write json", err)
	//	}
	//}
	for _, conn := range reg.GetConnections(ctx.SessionId) {
		err := conn.WriteMessage(websocket.TextMessage, jsonData)
		if err != nil {
			reg.logger.Error("WsHandler handleRead error to write json", "err", err,
				"userId", ctx.UserId,
				"userName", ctx.UserName,
				"data", jsonData,
			)
			continue
		}
	}
	//for _, conn := range reg.GetConnections(ctx.SessionId) {
	//	message := make(map[string]string)
	//	message[ctx.UserId] = ctx.UserName
	//	jsonData, err = json.Marshal(message)
	//	if err != nil {
	//		reg.logger.Error("WsHandler handleRead error to marshal json",
	//			"err", err)
	//		return
	//	}
	//	if conn != ctx.Conn {
	//		err = conn.WriteMessage(websocket.TextMessage, jsonData)
	//		if err != nil {
	//			reg.logger.Error("WsHandler handleRead error to write json", "err", err)
	//		}
	//		reg.logger.Info("ws sends to all", "message", message)
	//	}
	//}
}

func handleDelete(sessionID, userID string, reg *ConnectionRegistry) {
	reg.logger.Info("ws admin delete user user", "userId", userID, "SessionId", sessionID)
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
		err := conn.WriteMessage(websocket.TextMessage, jsonData)
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

func (r *ConnectionRegistry) HandleCloseConnection(userId, userName string, conn *websocket.Conn) {
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			r.logger.Info("Connection closed:", "userId", userId,
				"userName", userName,
				"err", err)
			break
		}
	}
}
