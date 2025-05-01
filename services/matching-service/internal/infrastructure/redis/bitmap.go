package redis

import (
	"fmt"

	"github.com/royroki/matching-service/internal/domain"
)

func getKey(module, tag string) string {
	return fmt.Sprintf("waiting:%s:%s", module, tag)
}

func AddUserToTags(req domain.MatchRequest) error {
	for _, tag := range req.Tags {
		key := getKey(req.Module, tag)
		err := Rdb.SetBit(Ctx, key, int64(req.UserID), 1).Err()
		if err != nil {
			return err
		}
	}
	return nil
}
