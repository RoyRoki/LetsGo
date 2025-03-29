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
	user, err := s.userRepo.GetUser(ctx, userID)
	if err != nil {
		log.Printf("Error retrieving user for user id %s\n", userID)
		return nil, err
	}

	chat, err := s.chatRepo.GetChatSession(ctx, user.ChatID)
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
	user, err := s.userRepo.GetUser(ctx, userID)
	chat, err := s.chatRepo.GetChatSession(ctx, user.ChatID)
	if err != nil {
		log.Printf("Error retrieving chat session for delete: %v", err)
		return err
	}

	// Send message for disconnect
	if chat.UserA.UserID != userID {
		err = s.wsRepo.SendMessage(chat.UserA.UserID, []byte("Wait for new partner..."))
		if err == nil {
			s.userRepo.AddUserToQueue(ctx, chat.UserA)
		}
	} else {
		err = s.wsRepo.SendMessage(chat.UserB.UserID, []byte("Wait for new partner..."))
		if err == nil {
			s.userRepo.AddUserToQueue(ctx, chat.UserB)
		}
	}

	err = s.chatRepo.DeleteChatSession(ctx, chat.ID)
	if err != nil {
		log.Printf("Error deleting chat session: %v", err)
		return err
	}
	s.wsRepo.RemoveConnection(userID)
	s.userRepo.RemoveUser(ctx, userID)

	log.Printf("Chat session")
	return nil
}

// AddUserToQueue calls UserRepository method via ChatService
func (s *ChatService) AddUserToQueue(ctx context.Context, user entity.User) error {
	err := s.userRepo.AddUserToQueue(ctx, user)
	if err == nil {
		message := fmt.Sprintf("ID : %s \nLetsGo wait for partner....", user.UserID)
		s.wsRepo.SendMessage(user.UserID, []byte(message))
	}
	return err
}

// CreateChatSession saves a chat session in Redis
func (s *ChatService) CreateChatSession(ctx context.Context, chat *entity.Chat) error {
	err := s.chatRepo.SaveChatSession(ctx, chat)
	if err != nil {
		log.Printf("❌ Error saving chat session: %v", err)
		return err
	}
	// Notify users about their new chat partner
	s.chatRepo.NotifyPartnerUpdate(ctx, chat.UserA.UserID, &chat.UserB)
	s.chatRepo.NotifyPartnerUpdate(ctx, chat.UserB.UserID, &chat.UserA)

	// Send message for new chat
	s.wsRepo.SendMessage(chat.UserA.UserID, []byte("💬 Open new chat with "+string(chat.UserB.UserID)))
	s.wsRepo.SendMessage(chat.UserB.UserID, []byte("💬 Open new chat with "+string(chat.UserA.UserID)))

	return nil
}

// ListenFromConnection listens for messages from a connected user
func (s *ChatService) ListenFromConnection(userID string) {
	ctx := context.Background()

	ws := s.wsRepo.GetConnection(userID)
	if ws == nil {
		log.Printf("⚠️ No active WebSocket connection for user %s", userID)
		return
	}

	user, err := s.userRepo.GetUser(ctx, userID)
	if err != nil {
		log.Println("Failed to listen.")
	}

	partnerUpdates := s.chatRepo.SubscribeToChatUpdates(ctx, userID)
	partner, err := s.chatRepo.GetChatPartner(ctx, user.ChatID, user.UserID)
	if err != nil {
		log.Println("⚠️ Failed to fetch initial chat partner")
		return
	}

	// Goroutine to handle partner updates
	go func() {
		for newPartner := range partnerUpdates {
			partner = newPartner // Update partner dynamically
			message := "🔄 Your chat partner has changed."
			s.wsRepo.SendMessage(userID, []byte(message))
		}
	}()

	defer func() {
		message := "Your partner is disconnected."
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
			s.wsRepo.SendMessage(partner.UserID, []byte("Partner Connection failed!"))
			log.Printf("⚠️ Error reading message from %s: %v", userID, err)
			break // Exit loop on error (disconnect)
		}

		log.Printf("📩 Received message from %s: %s", userID, string(message))

		// Forward the message to the user's chat partner
		err = s.wsRepo.SendMessage(partner.UserID, message)
		if err != nil {
			s.wsRepo.SendMessage(userID, []byte("Server Failed!"))
			log.Printf("⚠️ Error forwarding message: %v", err)
			break
		}
	}
}

// UpdateUserChatID updates the user's chat ID
func (s *ChatService) UpdateUserChatID(ctx context.Context, userID, chatID string) error {
	return s.userRepo.UpdateUserChatID(ctx, userID, chatID)
}
