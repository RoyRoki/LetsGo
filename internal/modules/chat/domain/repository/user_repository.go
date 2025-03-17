package repository

import (
	"context"

	"github.com/royroki/LetsGo/internal/modules/chat/domain/entity"
)

type UserRepository interface {
	AddUserToQueue(ctx context.Context, user entity.User) error
	GetUser(ctx context.Context, userID string) (*entity.User, error)
	PopTopUsers(ctx context.Context, i int) ([]entity.User, error)
	RemoveUser(ctx context.Context, userID string) error
}
