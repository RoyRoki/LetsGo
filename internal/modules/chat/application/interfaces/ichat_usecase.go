package interfaces

import (
	"context"

	"github.com/royroki/LetsGo/internal/modules/chat/domain/entity"
)

// ChatUseCase defines the use case contract
type ChatUseCase interface {
	GetChatPartner(ctx context.Context, userID string) (any, error)
	EndChatSession(ctx context.Context, userID string) error
	HandleNewConnection(ctx context.Context, connId, userId string) error
	HandleChatPair(ctx context.Context, userA, userB entity.User) error
	ListenFromConnection(userID string)
}
