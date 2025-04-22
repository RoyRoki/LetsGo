package service

import (
	"fmt"

	enum "github.com/royroki/LetsGo/internal/common/matching/domain/enums"
	"github.com/royroki/LetsGo/internal/common/matching/domain/repository"
)

type MatchingService struct {
	BitMapRepo repository.BitMapRepository
}

func NewMatchingService(bitmapRepo repository.BitMapRepository) *MatchingService {
	return &MatchingService{
		BitMapRepo: bitmapRepo,
	}
}

// findTagMatch checks if there is an exact match for the user based on their tags.
func (s *MatchingService) findTagMatch(module enum.Module, tags []string) (uint32, error) {
	for _, tag := range tags {
		// Retrive the bitmap from the given tag
		bitmap, err := s.BitMapRepo.GetBitMap(module, tag)
		if err != nil {
			return 0, fmt.Errorf("error retrieving bitmap for tag %s: %v", tag, err)
		}
		iter := bitmap.Iterator()

		if iter.HasNext() {
			userId := iter.Next()
			return userId, nil
		}
	}
	return 0, nil
}

// MatchUser attempts to find a user to match with the provided tags.
// If an exact match is not found within 5 seconds, it will attempt to match based on FIFO.
func (s *MatchingService) MatchUser(module enum.Module, tags []string, userID uint32) (uint32, bool) {
	var matchedUserID uint32
	var err error

	for i := 0; i < 5; i++ {
		matchedUserID, err = s.findTagMatch(module, tags)
		if err != nil {
			fmt.Printf("Error during tag matching: %v\n", err)
		}

	}
}
