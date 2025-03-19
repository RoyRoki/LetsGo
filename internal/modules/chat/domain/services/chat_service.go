package service

import (
	"context"
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
	return s.userRepo.AddUserToQueue(ctx, user)
}

// CreateChatSession saves a chat session in Redis
func (s *ChatService) CreateChatSession(ctx context.Context, chat *entity.Chat) error {
	return s.chatRepo.SaveChatSession(ctx, chat)
}

// ForwardMessage forwards a message to the chat partner
func (s *ChatService) ForwardMessage(senderID string, message []byte) error {
	ctx := context.Background()

	// Retrieve chat session
	chat, err := s.chatRepo.GetChatSession(ctx, senderID)
	if err != nil {
		log.Printf("⚠️ Error retrieving chat session for user %s: %v", senderID, err)
		return err
	}

	// Determine recipient
	recipientID := chat.UserA.UserID
	if chat.UserA.UserID == senderID {
		recipientID = chat.UserB.UserID
	}

	// Use WebSocketRepository instead of direct WebSocketHub access
	err = s.wsRepo.SendMessage(recipientID, message)
	if err != nil {
		log.Printf("⚠️ Error sending message to %s: %v", recipientID, err)
	}
	return err
}
