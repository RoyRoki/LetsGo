package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/royroki/LetsGo/internal/modules/chat/domain/entity"
	"github.com/royroki/LetsGo/internal/modules/chat/domain/repository"
)

// RedisChatRepository handles chat session storage in Redis
type ChatRepository struct {
	client *redis.Client
}

// NewRedisChatRepository initializes a new RedisChatRepository
func NewChatRepository(client *redis.Client) repository.ChatRepository {
	return &ChatRepository{client: client}
}

// SaveChatSession stores a chat session in Redis
func (r *ChatRepository) SaveChatSession(ctx context.Context, chat *entity.Chat) error {
	chatKey := fmt.Sprintf("chat:%s", chat.ID)

	// Convert chat struct to JSON
	chatData, err := json.Marshal(chat)
	if err != nil {
		log.Printf("Error marshalling chat data: %v", err)
		return err
	}

	_, err = r.client.HSet(ctx, chatKey, map[string]interface{}{
		"userA":     chat.UserA.UserID,
		"userB":     chat.UserB.UserID,
		"startTime": chat.StartTime.Unix(),
		"endTime":   0, // 0 means chat is ongoing
		"data":      string(chatData),
	}).Result()

	if err != nil {
		log.Printf("Error storing chat session: %v", err)
		return err
	}
	log.Printf("Chat session started: %s <-> %s", chat.UserA.UserID, chat.UserB.UserID)
	return nil
}

// GetChatSession retrieves a chat session from Redis
func (r *ChatRepository) GetChatSession(ctx context.Context, chatID string) (*entity.Chat, error) {
	chatKey := fmt.Sprintf("chat:%s", chatID)

	// Retrieve chat data from Redis
	data, err := r.client.HGetAll(ctx, chatKey).Result()
	if err != nil || len(data) == 0 {
		return nil, fmt.Errorf("chat session not found: %s", chatID)
	}

	// Convert JSON back to struct
	var chat entity.Chat
	err = json.Unmarshal([]byte(data["data"]), &chat)
	if err != nil {
		return nil, fmt.Errorf("error decoding chat session: %v", err)
	}

	return &chat, nil
}

// GetChatPartner retrieves the chat partner for a user
func (r *ChatRepository) GetChatPartner(ctx context.Context, chatID, userID string) (*entity.User, error) {
	// Retrieve the chat session
	chat, err := r.GetChatSession(ctx, chatID)
	if err != nil {
		return nil, err
	}

	// Return the chat partner
	if chat.UserA.UserID == userID {
		return &chat.UserB, nil
	}
	return &chat.UserA, nil
}

// DeleteChatSession removes a chat session from Redis
func (r *ChatRepository) DeleteChatSession(ctx context.Context, chatID string) error {
	chatKey := fmt.Sprintf("chat:%s", chatID)

		// Delete chat session
	err := r.client.Del(ctx, chatKey).Err()
	if err != nil {
		log.Printf("Error deleting chat session: %v", err)
		return err
	}
	log.Printf("Chat session deleted: %s", chatID)
	return nil
}



// SubscribeToChatUpdates listens for partner changes in Redis.
func (r *ChatRepository) SubscribeToChatUpdates(ctx context.Context, userID string) <-chan *entity.User {
	channel := "chat_updates:" + userID
	sub := r.client.Subscribe(ctx, channel)

	updates := make(chan *entity.User, 1) // Buffered to prevent blocking

	go func() {
		defer sub.Close()
		for msg := range sub.Channel() {
			partner := &entity.User{UserID: msg.Payload}
			updates <- partner
		}
		close(updates) // Close when subscription ends
	}()

	return updates
}

// NotifyPartnerUpdate publishes a new chat partner to Redis.
func (r *ChatRepository) NotifyPartnerUpdate(ctx context.Context, userID string, partner *entity.User) {
	channel := "chat_updates:" + userID
	err := r.client.Publish(ctx, channel, partner.UserID).Err()
	if err != nil {
		log.Printf("❌ Redis publish failed for %s: %v", userID, err)
	}
}
