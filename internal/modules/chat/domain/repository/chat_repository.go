package repository

import (
	"context"

	"github.com/royroki/LetsGo/internal/modules/chat/domain/entity"
)

// ChatRepository defines the contract for managing chat sessions.
type ChatRepository interface {
	// Save a new chat session
	SaveChatSession(ctx context.Context, chat *entity.Chat) error

	// Retrieve a chat session by ID
	GetChatSession(ctx context.Context, chatID string) (*entity.Chat, error)

	// Get the chat partner for a user
	GetChatPartner(ctx context.Context, chatID, userID string) (*entity.User, error)

	// Delete a chat session from storage
	DeleteChatSession(ctx context.Context, chatID string) error

	// Subcribe for chat updates
	SubscribeToChatUpdates(ctx context.Context, userID string) <-chan *entity.User 

	// Notify the chat updates
	NotifyPartnerUpdate(ctx context.Context, userID string, partner *entity.User)
}
