package repository

import enum "github.com/royroki/LetsGo/internal/common/matching/domain/enums"

type QueueRepository interface {
	EnQueue(module enum.Module, userID uint32) error
	DeQueue(module enum.Module) (uint32, error)
	RemoveUser(module enum.Module, userID uint32) error
}
