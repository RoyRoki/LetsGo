package web_socket

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/royroki/LetsGo/internal/modules/chat/application/interfaces"
	web_socket "github.com/royroki/LetsGo/internal/modules/chat/infrastructure/websocket"
)

// WebSocketHub manages active WebSocket connections.
type WebSocketHandler struct {
	useCase  interfaces.ChatUseCase
	upgrader websocket.Upgrader
	wsHub    *web_socket.WebSocketHub
}

// NewWebSocketHub initializes WebSocketHub with the chat use case.
func NewWebSocketHandler(useCase interfaces.ChatUseCase, hub *web_socket.WebSocketHub) *WebSocketHandler {
	return &WebSocketHandler{
		useCase: useCase,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
		wsHub: hub,
	}
}

// HandleWSConnection upgrades the HTTP request to WebSocket and handles the connection lifecycle.
func (h *WebSocketHandler) HandleWSConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	// Generate connID and extract userID
	userID := uuid.New().String()

	// Add the new connection to ws hub
	h.wsHub.AddConnection(userID, conn)

	// Inform use case of new connection
	err = h.useCase.HandleNewConnection(r.Context(), userID)
	if err != nil {
		log.Printf("Error connecting user: %v", err)
		h.wsHub.RemoveConnection(userID)
		return
	}

}
