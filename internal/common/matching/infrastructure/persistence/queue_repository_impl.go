package persistence

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	enum "github.com/royroki/LetsGo/internal/common/matching/domain/enums"
	"github.com/royroki/LetsGo/internal/common/matching/domain/repository"
)

type QueueRepositoryImpl struct {
	Redis *redis.Client
}

func NewQueueRepository(rdb *redis.Client) repository.QueueRepository {
	return &QueueRepositoryImpl{Redis: rdb}
}

func (q *QueueRepositoryImpl) queueKey(module string) string {
	return fmt.Sprintf("%s:queue", module)
}

func (q *QueueRepositoryImpl) EnQueue(module enum.Module, userID uint32) error {
	ctx := context.Background()
	return q.Redis.RPush(ctx, q.queueKey(string(module)), userID).Err()
}

func (q *QueueRepositoryImpl) DeQueue(module enum.Module) (uint32, error) {
	ctx := context.Background()
	res, err := q.Redis.LPop(ctx, q.queueKey(string(module))).Uint64()
	return uint32(res), err
}

func (q *QueueRepositoryImpl) RemoveUser(module enum.Module, userID uint32) error {
	ctx := context.Background()
	return q.Redis.LRem(ctx, q.queueKey(string(module)), 0, userID).Err()
}
