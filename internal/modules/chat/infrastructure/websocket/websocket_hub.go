package web_socket_hub

import (
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/royroki/LetsGo/internal/modules/chat/domain/repository"
)

// WebSocketHub manages WebSocket connections
type WebSocketHub struct {
	WSHub map[string]*websocket.Conn
	mu    sync.Mutex // Protects concurrent access
}

// Ensure WebSocketHub implements WebSocketRepository
var _ repository.WebSocketRepository = &WebSocketHub{}

// NewWebSocketHub initializes WebSocketHub
func NewWebSocketHub() *WebSocketHub {
	return &WebSocketHub{WSHub: make(map[string]*websocket.Conn)}
}

// AddConnection stores a WebSocket connection
func (h *WebSocketHub) AddConnection(userID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.WSHub[userID] = conn
}

// RemoveConnection removes a WebSocket connection
func (h *WebSocketHub) RemoveConnection(userID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if conn, exists := h.WSHub[userID]; exists {
		conn.Close()
		delete(h.WSHub, userID) // Remove from hub
	}
}

// GetConnection retrieves a WebSocket connection
func (h *WebSocketHub) GetConnection(userID string) *websocket.Conn {
	h.mu.Lock()
	defer h.mu.Unlock()

	conn, exists := h.WSHub[userID]
	if !exists || conn == nil {
		log.Printf("⚠️ No active WebSocket connection for %s", userID)
		return nil
	}
	return conn
}

// SendMessage sends a message to a connected user
func (h *WebSocketHub) SendMessage(userID string, message []byte) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	conn, exists := h.WSHub[userID]
	if !exists {
		return fmt.Errorf("user %s not connected", userID)
	}

	err := conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		return fmt.Errorf("error sending message to %s: %v", userID, err)
	}
	return nil
}

// Shutdown gracefully closes all WebSocket connections and clears the hub
func (h *WebSocketHub) Shutdown() {
	h.mu.Lock()
	defer h.mu.Unlock()

	log.Println("🔻 Closing all active WebSocket connections...")

	for userID, conn := range h.WSHub {
		err := conn.Close()
		if err != nil {
			log.Printf("⚠️ Error closing WebSocket for user %s: %v", userID, err)
		}
		delete(h.WSHub, userID)
	}

	log.Println("✅ WebSocketHub shutdown complete.")
}

