package Handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"log"
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

		// Ensure session exists in registry
		registry.RegisterSession(token.SessionId)

		// Register this connection
		if err := registry.RegisterConnection(token.SessionId, token.UserId, token.UserName, conn); err != nil {
			fmt.Println("error onRegisterConnection", err, token.SessionId, token.UserId, conn)
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
		welcome := fmt.Sprintf(`{"type":"welcome","sessionId":"%s","userId":"%s"}`, ctx.SessionId, ctx.UserId)
		conn.WriteMessage(websocket.TextMessage, []byte(welcome))

		// Start reading messages for this connection in a separate goroutine.
		go handleRead(ctx, registry)
	}
}

//человек делает запрос на вебсокет с его токеном. Далее Всем юзерам в такой руме приходит
// сообщение, если только в этой руме уже небыло чела с таким же именем и юзеркодом

type ConnectionRegistry struct {
	mu          sync.RWMutex
	connections map[string]map[string]*websocket.Conn
	rooms       map[string]map[string]string
}

// NewConnectionRegistry initializes the ConnectionRegistry
func NewConnectionRegistry() *ConnectionRegistry {
	return &ConnectionRegistry{
		connections: make(map[string]map[string]*websocket.Conn),
		rooms:       make(map[string]map[string]string),
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
	return nil
}

// UnregisterConnection removes joined user connection, (e.g., on user disconnect)
func (r *ConnectionRegistry) UnregisterConnection(sessionID, userID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if sessions, exists := r.connections[sessionID]; exists {
		delete(sessions, userID)
	}
	if rooms, exists := r.rooms[sessionID]; exists {
		delete(rooms, userID)
	}
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

func (r *ConnectionRegistry) GetRooms(sessionID string) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var rooms []string
	if sessions, exists := r.rooms[sessionID]; exists {
		for _, c := range sessions {
			rooms = append(rooms, c)
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
		fmt.Println("Malformed token")
	default:
		fmt.Println("Couldn't handle this token:", err)
	}

	claims, ok := token.Claims.(*shared.UserToken)
	fmt.Println("Decoded token: ", *claims)
	if !ok {
		return nil, fmt.Errorf("error decoding jwt")
	}

	return claims, nil
}

// Всем юзерам он должен отправить свое имя, а себе должен отправить и свое имя и весь список
func handleRead(ctx *ConnectionContext, reg *ConnectionRegistry) {
	err := reg.RegisterConnection(ctx.SessionId, ctx.UserId, ctx.UserName, ctx.Conn)
	if err != nil {
		fmt.Println("ws error in token: %w", err)
		return
	}
	fmt.Println("ws connected to session", len(reg.GetConnections(ctx.SessionId)))
	m := reg.GetRooms(ctx.SessionId)
	con := reg.connections[ctx.SessionId][ctx.UserId]
	jsonData, err := json.Marshal(m)
	if err != nil {
		log.Println("Ошибка сериализации:", err)
		return
	}
	if len(m) > 0 {
		fmt.Println(m)
		err = con.WriteMessage(websocket.TextMessage, jsonData)
		if err != nil {
			fmt.Println("ws error in token: %w", err)
		}
	}
	for _, conn := range reg.GetConnections(ctx.SessionId) {
		jsonData, err = json.Marshal(ctx.UserId)
		if err != nil {
			log.Println("Ошибка сериализации:", err)
			return
		}
		if conn != ctx.Conn {
			err = conn.WriteMessage(websocket.TextMessage, jsonData)
			if err != nil {
				fmt.Println("ws error in token: %w", err)
			}
		}
	}
}
