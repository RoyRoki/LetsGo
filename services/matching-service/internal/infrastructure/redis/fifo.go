package redis

import (
	"fmt"
	"log"
	"strconv"

	redis "github.com/redis/go-redis/v9"
)

func FifoMatch(module string, userID int32) (partnerID int32, matched bool) {
	// Try to pop a partner from FIFO queue
	partnerStr, err := Rdb.RPop(Ctx, fmt.Sprintf("queue:%s:fifo", module)).Result()

	// If the queue is empty, push the current user and return false
	if err == redis.Nil {
		log.Printf("No partner found in queue for module %s. Adding user %d to the queue.", module, userID)
		Rdb.LPush(Ctx, fmt.Sprintf("queue:%s:fifo", module), strconv.Itoa(int(userID)))
		return 0, false
	} else if err != nil {
		log.Printf("Redis error: %v", err)
		return 0, false
	}

	// Convert partnerStr to int and check if it's the same as the current user
	partnerIDInt, err := strconv.Atoi(partnerStr)
	if err != nil {
		log.Printf("Error converting partner ID from string: %v", err)
		return 0, false
	}

	if int32(partnerIDInt) == userID {
		// If the user is matched with themselves, push them back and return false
		log.Printf("User %d is in the queue but was matched with themselves. Re-adding to the queue.", userID)
		Rdb.LPush(Ctx, fmt.Sprintf("queue:%s:fifo", module), strconv.Itoa(int(userID)))
		return 0, false
	}

	// Partner found and is not the current user
	log.Printf("User %d matched with partner %d in module %s", userID, partnerID, module)
	return int32(partnerIDInt), true
}
