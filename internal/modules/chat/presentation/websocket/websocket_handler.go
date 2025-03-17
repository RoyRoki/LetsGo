package websocket

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/royroki/LetsGo/internal/modules/chat/application/usecase"
)

// WebSocketHub manages active WebSocket connections.
type WebSocketHub struct {
	useCase  *usecase.ChatUseCase
	upgrader websocket.Upgrader
}

// NewWebSocketHub initializes WebSocketHub with the chat use case.
func NewWebSocketHub(useCase *usecase.ChatUseCase) *WebSocketHub {
	return &WebSocketHub{
		useCase: useCase,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
	}
}

// HandleWSConnection upgrades the HTTP request to WebSocket and handles the connection lifecycle.
func (hub *WebSocketHub) HandleWSConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := hub.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	// Generate connID and extract userID
	userID := r.RemoteAddr
	connID := conn.RemoteAddr().String()

	// Inform use case of new connection
	err = hub.useCase.HandleWSConnection(r.Context(), connID, userID)
	if err != nil {
		log.Printf("Error connecting user: %v", err)
		conn.Close()
		return
	}

	// Start message listener
	go hub.listenForMessages(context.Background(), userID, conn)
}

// listenForMessages listens for incoming messages and forwards them to paired user.
func (hub *WebSocketHub) listenForMessages(ctx context.Context, userID string, conn *websocket.Conn) {
	defer func() {
		conn.Close()
		delete(hub.useCase.Connections, userID)
		hub.useCase.EndChatSession(ctx, userID)
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Connection closed for user %s: %v", userID, err)
			break
		}

		// Retrieve chat partner via use case
		partner, err := hub.useCase.GetChatPartner(ctx, userID)
		if err != nil {
			log.Printf("Partner not found for user %s: %v", userID, err)
			continue
		}

	}
}
