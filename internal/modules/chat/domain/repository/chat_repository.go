package repository

import "github.com/royroki/LetsGo/internal/modules/chat/domain/entity"

type ChatRepository interface {
	CreateChat(chat *entity.Chat) error
	GetChat(id string) (*entity.Chat, error)
	UpdateChat(chat *entity.Chat) error
	DeleteChat(id string) error
}
