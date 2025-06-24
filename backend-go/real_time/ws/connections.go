package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
)

// ConnectionRegistry manages all users' ws connections
type ConnectionRegistry struct {
	mu          sync.RWMutex
	connections map[string]map[string]*websocket.Conn
}

// NewConnectionRegistry initializes the ConnectionRegistry
func NewConnectionRegistry() *ConnectionRegistry {
	return &ConnectionRegistry{
		connections: make(map[string]map[string]*websocket.Conn),
	}
}

// RegisterSession creates a new session entry; idempotent
func (r *ConnectionRegistry) RegisterSession(sessionID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.connections[sessionID]; !exists {
		r.connections[sessionID] = make(map[string]*websocket.Conn)
	}
}

// UnregisterSession removes session entirely (e.g., on session end)
func (r *ConnectionRegistry) UnregisterSession(sessionID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.connections, sessionID)
}

// RegisterConnection adds new joined user connection, mapping to a corresponding session
func (r *ConnectionRegistry) RegisterConnection(sessionID, userID string, conn *websocket.Conn) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	sessions, exists := r.connections[sessionID]
	if !exists {
		return fmt.Errorf("session %s not found", sessionID)
	}
	sessions[userID] = conn
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
