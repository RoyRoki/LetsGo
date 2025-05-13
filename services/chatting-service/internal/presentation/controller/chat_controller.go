package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/royroki/services/chatting-service/internal/app/models"
	"github.com/royroki/services/chatting-service/internal/app/usecases"
	"github.com/royroki/services/chatting-service/internal/constants"
	"github.com/royroki/services/chatting-service/internal/domain/entity"
	"github.com/royroki/services/chatting-service/internal/infrastructure/websocket"
)

type ChatController struct {
	wsServer    *websocket.WCServer
	chatUsecase *usecases.ChatUsecase
}

func NewChatController(ws *websocket.WCServer, uc *usecases.ChatUsecase) *ChatController {
	return &ChatController{wsServer: ws, chatUsecase: uc}
}

// HandleWebSocket upgrades the connection and registers a new user
func (c *ChatController) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade to WebSocket
	conn, userID, err := c.wsServer.Upgrade(w, r)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
		return
	}

	// Create user entity
	user := &entity.User{
		ID:     userID,
		Conn:   conn,
		IP:     r.RemoteAddr,
		Module: "chat",
	}

	// Register user
	c.wsServer.Register(user)

	// Initial success message
	conn.WriteJSON(entity.Message{
		From:    0,
		Content: "Connected. Send tags to continue.",
		Status:  constants.ServerConnected,
	})

	// Listen for messages
	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			log.Printf("User %d disconnected or error: %v\n", userID, err)
			c.wsServer.Unregister(userID)
			break
		}

		var chatReq models.ChatRequest
		if err := json.Unmarshal(msgBytes, &chatReq); err != nil {
			log.Printf("Invalid request from user %d: %v\n", userID, err)
			conn.WriteJSON(entity.Message{
				From:    0,
				Content: "Invalid request format.",
				Status:  constants.InvalidRequestFormat,
			})
			continue
		}

		// Store tags/module in the user entity
		user.Tags = chatReq.Tags

		log.Printf("User %d tags: %v, module: %s", userID, user.Tags, user.Module)

		// Handle matching logic
		go c.chatUsecase.HandleMatchRequest(user)
	}
}
