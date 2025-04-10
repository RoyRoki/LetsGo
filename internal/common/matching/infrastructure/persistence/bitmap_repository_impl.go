package persistence

import (
	"bytes"
	"context"
	"fmt"

	"github.com/RoaringBitmap/roaring"
	"github.com/redis/go-redis/v9"
	enum "github.com/royroki/LetsGo/internal/common/matching/domain/enums"
	"github.com/royroki/LetsGo/internal/common/matching/domain/repository"
)

type BitMapRepositoryImpl struct {
	Redis *redis.Client
}

func NewBitMapRepository(rdb *redis.Client) repository.BitMapRepository {
	return &BitMapRepositoryImpl{Redis: rdb}
}

func (r *BitMapRepositoryImpl) bitmapKey(module, tag string) string {
	return fmt.Sprintf("%s:tag:%s", module, tag)
}

func (r *BitMapRepositoryImpl) GetBitMap(module enum.Module, tag string) (*roaring.Bitmap, error) {
	ctx := context.Background()
	data, err := r.Redis.Get(ctx, r.bitmapKey(string(module), tag)).Bytes()
	if err == redis.Nil {
		return roaring.New(), nil
	} else if err != nil {
		return nil, err
	}
	bitmap := roaring.New()
	err = bitmap.UnmarshalBinary(data)
	return bitmap, err
}

func (r *BitMapRepositoryImpl) AddToBitMap(module enum.Module, tag string, userID uint32) error {
	bitmap, err := r.GetBitMap(module, tag)
	if err != nil {
		return err
	}
	bitmap.Add(userID)
	var buf bytes.Buffer
	if _, err = bitmap.WriteTo(&buf); err != nil {
		return err
	}
	ctx := context.Background()
	return r.Redis.Set(ctx, r.bitmapKey(string(module), tag), buf.Bytes(), 0).Err()
}

func (r *BitMapRepositoryImpl) RemoveFromBitMap(module enum.Module, tag string, userID uint32) error {
	bitmap, err := r.GetBitMap(module, tag)
	if err != nil {
		return err
	}
	bitmap.Remove(userID)
	var buf bytes.Buffer
	if _, err := bitmap.WriteTo(&buf); err != nil {
		return err
	}
	ctx := context.Background()
	return r.Redis.Set(ctx, r.bitmapKey(string(module), tag), buf.Bytes(), 0).Err()
}
