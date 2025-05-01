package lua

import "github.com/royroki/matching-service/internal/domain"

func MatchUser(req domain.MatchRequest) (domain.MatchResult, error) {
	// Placeholder: In final version this will call a Lua script.
	return domain.MatchResult{Matched: false}, nil
}
