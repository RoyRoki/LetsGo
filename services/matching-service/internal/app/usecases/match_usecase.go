package usecases

import (
	"github.com/royroki/matching-service/internal/domain"
	"github.com/royroki/matching-service/internal/infrastructure/lua"
	"github.com/royroki/matching-service/internal/infrastructure/redis"
)

type MatchUsercase struct{}

func NewMatchUsecase() *MatchUsercase {
	return &MatchUsercase{}
}

// TryMatch attempts to find a partner for the given user and module.
func (u *MatchUsercase) TryMatch(req domain.MatchRequest) (domain.MatchResult, error) {
	// First Try to find a match
	result, err := lua.MatchUser(req)
	if err != nil {
		return domain.MatchResult{}, err
	}

	// If no match found, add user to bitmap queues
	if !result.Matched {
		if err := redis.AddUserToTags(req); err != nil {
			return domain.MatchResult{}, err
		}
	}

	return result, nil
}
