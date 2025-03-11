package datasource

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
	"github.com/royroki/LetsGo/internal/modules/chat/domain/entity"
	"github.com/royroki/LetsGo/internal/modules/chat/domain/repository"
)

type RedisChatRepository struct {
	client *redis.Client
}

func NewRedisChatRepository(client *redis.Client) repository.ChatRepository {
	return &RedisChatRepository{client: client}
}

func (r *RedisChatRepository) CreateChat(chat *entity.Chat) error {
	// Convert the chat object to JSON
	chatData, err := json.Marshal(chat)
	if err != nil {
		return err
	}

	// Save the chat in Redis (use the chat ID as the key)
	err = r.client.Set(context.Background(), chat.ID, chatData, 0).Err()
	return err
}

func (r *RedisChatRepository) GetChat(id string) (*entity.Chat, error) {
	chatData, err := r.client.Get(context.Background(), id).Result()
	if err != nil {
		return nil, err
	}

	var chat entity.Chat
	err = json.Unmarshal([]byte(chatData), &chat)
	if err != nil {
		return nil, err
	}

	return &chat, nil
}

func (r *RedisChatRepository) UpdateChat(chat *entity.Chat) error {
	chatData, err := json.Marshal(chat)
	if err != nil {
		return err
	}

	err = r.client.Set(context.Background(), chat.ID, chatData, 0).Err()
	return err
}

func (r *RedisChatRepository) DeleteChat(id string) error {
	return r.client.Del(context.Background(), id).Err()
}
