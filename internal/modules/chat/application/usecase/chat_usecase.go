package usecase

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/royroki/LetsGo/internal/modules/chat/domain/entity"
	"github.com/royroki/LetsGo/internal/modules/chat/domain/repository"
)

type ChatUseCase struct {
	chatRepo repository.ChatRepository
}

func NewChatUseCase(chatRepo repository.ChatRepository) *ChatUseCase {
	return &ChatUseCase{chatRepo: chatRepo}
}

func (c *ChatUseCase) CreateChat(participant1, participant2 string) (*entity.Chat, error) {
	// Generate a unique chat ID (you can use UUID)
	chatID := generateUniqueChatID()

	chat := &entity.Chat{
		ID:           chatID,
		Participant1: participant1,
		Participant2: participant2,
		Status:       "Active",         // You can start with an active status
		StartTime:    getCurrentTime(), // Timestamp when the chat starts
	}

	err := c.chatRepo.CreateChat(chat)
	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (c *ChatUseCase) GetChat(id string) (*entity.Chat, error) {
	chat, err := c.chatRepo.GetChat(id)
	if err != nil {
		return nil, err
	}

	if chat == nil {
		return nil, errors.New("chat not found")
	}

	return chat, nil
}

func (c *ChatUseCase) EndChat(id string) error {
	chat, err := c.chatRepo.GetChat(id)
	if err != nil {
		return err
	}

	if chat == nil {
		return errors.New("chat not found")
	}

	chat.Status = "Ended"           // Mark the chat as ended
	chat.EndTime = getCurrentTime() // Set the end time

	return c.chatRepo.UpdateChat(chat)
}

func generateUniqueChatID() string {
	return uuid.New().String()
}

func getCurrentTime() int64 {
	return time.Now().Unix()
}
