package repository

import "github.com/gorilla/websocket"

// WebSocketRepository defines WebSocket operations
type WebSocketRepository interface {
	AddConnection(userID string, conn *websocket.Conn)
	RemoveConnection(userID string)
	GetConnection(userID string) *websocket.Conn
	SendMessage(userID string, message []byte) error
	Shutdown()
}
