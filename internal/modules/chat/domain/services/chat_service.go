package service

import (
	"context"
	"fmt"
	"log"

	"github.com/royroki/LetsGo/internal/modules/chat/domain/entity"
	"github.com/royroki/LetsGo/internal/modules/chat/domain/repository"
)

// ChatService handles domain logic for chat
type ChatService struct {
	chatRepo repository.ChatRepository
	userRepo repository.UserRepository
	wsRepo   repository.WebSocketRepository
}

// NewChatService initializes ChatService
func NewChatService(chatRepo repository.ChatRepository, userRepo repository.UserRepository, wsRepo repository.WebSocketRepository) *ChatService {
	return &ChatService{chatRepo: chatRepo, userRepo: userRepo, wsRepo: wsRepo}
}

// GetChatPartner retrieves the chat partner of a user
func (s *ChatService) GetChatPartner(ctx context.Context, userID string) (*entity.User, error) {
	chat, err := s.chatRepo.GetChatSession(ctx, userID)
	if err != nil {
		log.Printf("Error retrieving chat session for user %s: %v", userID, err)
		return nil, err
	}

	if chat.UserA.UserID == userID {
		return &chat.UserB, nil
	}
	return &chat.UserA, nil
}

// EndChatSession removes a chat session and re-adds users to the queue
func (s *ChatService) EndChatSession(ctx context.Context, userID string) error {
	chat, err := s.chatRepo.GetChatSession(ctx, userID)
	if err != nil {
		log.Printf("Error retrieving chat session: %v", err)
		return err
	}

	// Send message for disconnect 
	s.wsRepo.SendMessage(chat.UserA.UserID, []byte("New the chat."))
	s.wsRepo.SendMessage(chat.UserB.UserID, []byte("New the chat."))

	// Remove chat session
	err = s.chatRepo.DeleteChatSession(ctx, chat.ID)
	if err != nil {
		log.Printf("Error deleting chat session: %v", err)
		return err
	}

	// Re-add users to waiting queue
	s.userRepo.AddUserToQueue(ctx, chat.UserA)
	s.userRepo.AddUserToQueue(ctx, chat.UserB)

	log.Printf("Chat session %s ended. Users %s and %s are back in queue.", chat.ID, chat.UserA.UserID, chat.UserB.UserID)
	return nil
}

// AddUserToQueue calls UserRepository method via ChatService
func (s *ChatService) AddUserToQueue(ctx context.Context, user entity.User) error {
	err := s.userRepo.AddUserToQueue(ctx, user)
	if err == nil {
		message := fmt.Sprintf("ID : %s \nLetsGo wait for partner", user.UserID)
		s.wsRepo.SendMessage(user.UserID, []byte(message))
	}
	return err;
}

// CreateChatSession saves a chat session in Redis
func (s *ChatService) CreateChatSession(ctx context.Context, chat *entity.Chat) error {
	err := s.chatRepo.SaveChatSession(ctx, chat)
	if err != nil {
		log.Printf("❌ Error saving chat session: %v", err)
		return err
	}
// Send message for new chat 
	s.wsRepo.SendMessage(chat.UserA.UserID, []byte("Ending the chat."))
	s.wsRepo.SendMessage(chat.UserB.UserID, []byte("New the chat."))


	log.Printf("✅ Chat session created: %s <-> %s (ChatID: %s)", chat.UserA.UserID, chat.UserB.UserID, chat.ID)
	return nil
}

// ForwardMessage forwards a message to the chat partner
func (s *ChatService) ForwardMessage(senderID string, message []byte) error {
	ctx := context.Background()

	user, err := s.userRepo.GetUser(ctx, senderID);

	// Retrieve chat session
	chat, err := s.chatRepo.GetChatSession(ctx, user.ChatID)
	if err != nil {
		log.Printf("⚠️ Error retrieving chat session for user %s: %v", user.UserID, err)
		return err
	}

	// Determine recipient
	recipientID := chat.UserA.UserID
	if chat.UserA.UserID == user.ChatID {
		recipientID = chat.UserB.UserID
	}

	// Use WebSocketRepository instead of direct WebSocketHub access
	err = s.wsRepo.SendMessage(recipientID, message)
	if err != nil {
		log.Printf("⚠️ Error sending message to %s: %v", recipientID, err)
	}
	return err
}

// ListenFromConnection listens for messages from a connected user
func (s *ChatService) ListenFromConnection(userID string) {
	ctx := context.Background()

	ws := s.wsRepo.GetConnection(userID)
	if ws == nil {
		log.Printf("⚠️ No active WebSocket connection for user %s", userID)
		return // Prevent further execution
	}
	user, err := s.userRepo.GetUser(ctx, userID)
	if err != nil {
		log.Println("Failed to listen.")
	}
	partner, err := s.chatRepo.GetChatPartner(ctx, user.ChatID, user.UserID)

	defer func() {
		message := "You are Disconnected."
		err := s.wsRepo.SendMessage(partner.UserID, []byte(message))

		if err != nil {
			log.Println("Failed to send disconnec message.")
		}
		// When user disconnects, remove from WebSocket hub and queue
		ws.Close()
		s.EndChatSession(context.Background(), userID)
		log.Printf("User disconnected: %s", userID)
	}()

	for {
		// Read incoming message
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Printf("⚠️ Error reading message from %s: %v", userID, err)
			break // Exit loop on error (disconnect)
		}

		log.Printf("📩 Received message from %s: %s", userID, string(message))

		// ✅ Forward the message to the user's chat partner
		err = s.ForwardMessage(partner.UserID, message)
		if err != nil {
			log.Printf("⚠️ Error forwarding message: %v", err)
			break
		}
	}
}

// UpdateUserChatID updates the user's chat ID
func (s *ChatService) UpdateUserChatID(ctx context.Context, userID, chatID string) error {
	return s.userRepo.UpdateUserChatID(ctx, userID, chatID)
}
