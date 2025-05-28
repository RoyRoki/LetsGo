package lua

import (
	_ "embed"
	"fmt"

	redis "github.com/redis/go-redis/v9"                      // add alias for clarity
	"github.com/royroki/letsgo/services/matching-service/internal/domain"
	Redis "github.com/royroki/letsgo/services/matching-service/internal/infrastructure/redis"
)

//go:embed match_user.lua
var matchUserScript string

var matchUserRedisScript = redis.NewScript(matchUserScript)

func MatchUser(req domain.MatchRequest) (domain.MatchResult, error) {
	// Prepare Redis KEYS (bitmap keys for each tag)
	keys := make([]string, len(req.Tags))
	for i, tag := range req.Tags {
		keys[i] = fmt.Sprintf("waiting:%s:%s", req.Module, tag)
	}

	// Cap scan to 4.29 billion (max 32-bit userID)
	const maxScan = 4294970000

	// Run Lua script
	res, err := matchUserRedisScript.Run(Redis.Ctx, Redis.Rdb, keys, req.UserID, maxScan).Int()
	if err != nil {
		return domain.MatchResult{Matched: false}, err
	}

	if res == -1 {
		return domain.MatchResult{Matched: false}, nil
	}

	return domain.MatchResult{
		Matched:   true,
		PartnerID: int32(res),
	}, nil
}
