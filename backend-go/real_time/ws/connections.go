package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"xxx/shared"
)

// ConnectionRegistry manages all users' ws connections
type ConnectionRegistry struct {
	mu          sync.RWMutex
	connections map[string]map[string]ConnectionContext // sessionId -> userId -> ConnectionContext
}

// NewConnectionRegistry initializes the ConnectionRegistry
func NewConnectionRegistry() *ConnectionRegistry {
	return &ConnectionRegistry{
		connections: make(map[string]map[string]ConnectionContext),
	}
}

// RegisterSession creates a new session entry; idempotent
func (r *ConnectionRegistry) RegisterSession(sessionID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.connections[sessionID]; !exists {
		r.connections[sessionID] = make(map[string]ConnectionContext)
		fmt.Println("Register new session:", r.connections)
	}
}

// UnregisterSession removes session entirely (e.g., on session end)
func (r *ConnectionRegistry) UnregisterSession(sessionID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.connections, sessionID)
}

// RegisterConnection adds new joined user connection, mapping to a corresponding session
func (r *ConnectionRegistry) RegisterConnection(ctx ConnectionContext) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.connections[ctx.SessionId]
	if !exists {
		return fmt.Errorf("session %s not found", ctx.SessionId)
	}
	r.connections[ctx.SessionId][ctx.UserId] = ctx
	fmt.Println("Register new connection:", r.connections)
	return nil
}

// UnregisterConnection removes joined user connection, (e.g., on user disconnect)
func (r *ConnectionRegistry) UnregisterConnection(sessionID, userID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if sessions, exists := r.connections[sessionID]; exists {
		delete(sessions, userID)
	}
}

// GetConnections gets a snapshot copy of connections to avoid holding lock during WriteMessage
func (r *ConnectionRegistry) GetConnections(sessionID string) []ConnectionContext {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var conns []ConnectionContext
	if sessions, exists := r.connections[sessionID]; exists {
		for _, ctx := range sessions {
			conns = append(conns, ctx)
		}
	}
	return conns
}

// BroadcastToSession sends the given payload to all users connected to a specific session.
// It uses the session ID to retrieve all active WebSocket connections and broadcasts the message.
func (r *ConnectionRegistry) BroadcastToSession(sessionId string, payload []byte, sendToAdmin bool) {
	receivers := r.GetConnections(sessionId)
	if !sendToAdmin {
		for i, rcv := range receivers {
			if rcv.Role == shared.RoleAdmin {
				receivers = append(receivers[:i], receivers[i+1:]...)
				fmt.Println("Receivers after admin removing:", receivers)
			}
		}
	}
	r.sendMessage(payload, receivers...)
}

// SendToAdmin sends the given payload only to the admin user of a specific session.
// It iterates over all connections in the session and filters by admin role before sending.
func (r *ConnectionRegistry) SendToAdmin(sessionId string, payload []byte) {
	fmt.Println("Sending to admin of", sessionId)
	fmt.Println(r.connections)
	for _, ctx := range r.connections[sessionId] {
		if ctx.Role == shared.RoleAdmin {
			fmt.Println("admin id: ", ctx.UserId)
			r.sendMessage(payload, ctx)
		}
	}
}

// sendMessage sends a WebSocket message (payload) to one or more connections.
// It logs errors but does not halt on failure to individual connections.
func (r *ConnectionRegistry) sendMessage(payload []byte, receivers ...ConnectionContext) {
	for _, ctx := range receivers {
		if ctx.Conn == nil {
			log.Println("Skipped sending message: connection is nil")
			continue
		}
		err := ctx.Conn.WriteMessage(websocket.TextMessage, payload)
		if err != nil {
			log.Printf("Failed to send message to connection: %v", err)
			r.UnregisterConnection(ctx.SessionId, ctx.UserId)
			continue
		}
		log.Println("Message sent successfully")
	}
}
