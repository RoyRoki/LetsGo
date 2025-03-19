package usecase

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/royroki/LetsGo/internal/modules/chat/application/interfaces"
	"github.com/royroki/LetsGo/internal/modules/chat/domain/entity"
	service "github.com/royroki/LetsGo/internal/modules/chat/domain/services"
)

type ChatUseCase struct {
	chatService *service.ChatService
}

// Ensure `ChatUseCaseImpl` implements `ChatUseCase`
var _ interfaces.ChatUseCase = &ChatUseCase{}

func NewChatUseCase(chatService *service.ChatService) *ChatUseCase {
	return &ChatUseCase{chatService: chatService}
}

func (c *ChatUseCase) GetChatPartner(ctx context.Context, userID string) (any, error) {
	return c.chatService.GetChatPartner(ctx, userID)
}

func (c *ChatUseCase) EndChatSession(ctx context.Context, userID string) error {
	return c.chatService.EndChatSession(ctx, userID)
}

// HandleWSConnection manages WebSocket connections, pairing users, and messaging.
func (c *ChatUseCase) HandleNewConnection(ctx context.Context, connId, userId string) error {
	log.Printf("User connected: %s (ConnID: %s)", userId, connId)

	// Create a User entity
	user := entity.User{
		UserID:   userId,
		ConnID:   connId,
		JoinTime: time.Now(),
		Chatted:  0,
	}

	// Add user to queue (Worker will pair them)
	err := c.chatService.AddUserToQueue(ctx, user)
	if err != nil {
		log.Printf("Error adding user to queue: %v", err)
		return err
	}

	log.Printf("User %s added to waiting queue. Worker will handle pairing.", userId)
	return nil
}

// HandleChatPair creates a chat session when two users are matched
func (c *ChatUseCase) HandleChatPair(ctx context.Context, userA, userB entity.User) error {
	// Create chat session entity
	chat := entity.Chat{
		ID:        uuid.New().String(),
		UserA:     userA,
		UserB:     userB,
		StartTime: time.Now(),
	}

	// Save chat session
	err := c.chatService.CreateChatSession(ctx, &chat)
	if err != nil {
		log.Printf("Error saving chat session: %v", err)
		return err
	}

	log.Printf("✅ Chat session started: %s <-> %s (ChatID: %s)", userA.UserID, userB.UserID, chat.ID)
	return nil
}

// ListenFromConnection listens for messages from a connected user
func (c *ChatUseCase) ListenFromConnection(userID string, conn *websocket.Conn) {
	defer func() {
		// When user disconnects, remove from WebSocket hub and queue
		conn.Close()
		c.chatService.EndChatSession(context.Background(), userID)
		log.Printf("User disconnected: %s", userID)
	}()

	for {
		// Read incoming message
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("⚠️ Error reading message from %s: %v", userID, err)
			break // Exit loop on error (disconnect)
		}

		log.Printf("📩 Received message from %s: %s", userID, string(message))

		// ✅ Forward the message to the user's chat partner
		err = c.chatService.ForwardMessage(userID, message)
		if err != nil {
			log.Printf("⚠️ Error forwarding message: %v", err)
			break
		}
	}
}
