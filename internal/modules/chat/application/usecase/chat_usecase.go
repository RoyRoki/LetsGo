package usecase

import (
	"context"

	"github.com/gorilla/websocket"
	"github.com/royroki/LetsGo/internal/modules/chat/domain/repository"
)

type ChatUseCase struct {
	chatRepo    repository.ChatRepository
	Connections map[string]*websocket.Conn
}

func (c *ChatUseCase) GetChatPartner(ctx context.Context, userID string) (any, error) {
	panic("unimplemented")
}

func (c *ChatUseCase) EndChatSession(ctx context.Context, userID string) {
	panic("unimplemented")
}

func NewChatUseCase(chatRepo repository.ChatRepository) *ChatUseCase {
	return &ChatUseCase{chatRepo: chatRepo}
}

// HandleWSConnection manages WebSocket connections, pairing users, and messaging.
func (c *ChatUseCase) HandleWSConnection(ctx context.Context, connId, userId string) error {

}
