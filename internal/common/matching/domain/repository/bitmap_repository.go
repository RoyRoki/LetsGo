package repository

import (
	"github.com/RoaringBitmap/roaring"
	enum "github.com/royroki/LetsGo/internal/common/matching/domain/enums"
)

type BitMapRepository interface {
	GetBitMap(module enum.Module, tag string) (*roaring.Bitmap, error)
	AddToBitMap(module enum.Module, tag string, userID uint32) error
	RemoveFromBitMap(module enum.Module, tag string, userID uint32) error
}
