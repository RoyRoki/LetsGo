package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/royroki/LetsGo/internal/modules/chat/domain/entity"
	"github.com/royroki/LetsGo/internal/modules/chat/domain/repository"
)

// UserRepository implements UserRepository using Redis
type UserRepository struct {
	client *redis.Client
	queue  string
}

// NewUserRepository initializes a Redis user repository
func NewUserRepository(client *redis.Client, queueName string) repository.UserRepository {
	return &UserRepository{
		client: client,
		queue:  queueName,
	}
}

// AddUserToQueue stores user entity in Redis and adds them to the queue
func (r *UserRepository) AddUserToQueue(ctx context.Context, user entity.User) error {
	priority := float64(time.Now().Unix()) // Lower score = higher priority

	// Store user in Redis Hash
	userKey := fmt.Sprintf("user:%s", user.UserID)
	userData, _ := json.Marshal(user)

	_, err := r.client.HSet(ctx, userKey, map[string]interface{}{
		"chatID":   user.ChatID,
		"joinTime": user.JoinTime.Unix(),
		"chatted":  user.Chatted,
		"data":     userData,
	}).Result()

	if err != nil {
		log.Printf("Error storing user entity: %v", err)
		return err
	}

	// Add user to the waiting queue (Sorted Set)
	_, err = r.client.ZAdd(ctx, r.queue, redis.Z{
		Score:  priority,
		Member: user.UserID,
	}).Result()

	if err != nil {
		log.Printf("Error adding user to queue: %v", err)
		return err
	}

	log.Printf("User %s added to the queue with priority %.0f", user.UserID, priority)
	return nil
}

// GetUser retrieves a user entity from Redis
func (r *UserRepository) GetUser(ctx context.Context, userID string) (*entity.User, error) {
	userKey := fmt.Sprintf("user:%s", userID)

	data, err := r.client.HGetAll(ctx, userKey).Result()
	if err != nil || len(data) == 0 {
		return nil, err // User not found
	}

	user := &entity.User{
		UserID:   userID,
		ChatID: data["chatID"],
		JoinTime: time.Unix(parseInt64(data["joinTime"]), 0),
		Chatted:  parseInt64(data["chatted"]),
	}

	return user, nil
}

// RemoveUser removes a user from Redis and the waiting queue
func (r *UserRepository) RemoveUser(ctx context.Context, userID string) error {
	userKey := fmt.Sprintf("user:%s", userID)

	// Remove user from Redis Hash
	err := r.client.Del(ctx, userKey).Err()
	if err != nil {
		return err
	}

	// Remove user from waiting queue (Sorted Set)
	err = r.client.ZRem(ctx, r.queue, userID).Err()
	return err
}

// Helper function to parse int64
func parseInt64(s string) int64 {
	val, err := time.ParseDuration(s + "s")
	if err != nil {
		return 0
	}
	return int64(val.Seconds())
}

// PopTopUsers retrieves and removes the top `limit` users from the queue
func (r *UserRepository) PopTopUsers(ctx context.Context, limit int) ([]entity.User, error) {
	// Step 1️⃣: Get the top `limit` user IDs from the queue
	userIDs, err := r.client.ZRange(ctx, r.queue, 0, int64(limit)-1).Result()
	if err != nil || len(userIDs) == 0 {
		log.Println("⚠️ No users found in queue")
		return nil, err
	}

	// Step 2️⃣: Remove these users from the queue
	r.client.ZRem(ctx, r.queue, userIDs)

	// Step 3️⃣: Retrieve full user data
	var users []entity.User
	for _, userID := range userIDs {
		user, err := r.GetUser(ctx, userID)
		if err == nil {
			users = append(users, *user)
		} else {
			log.Printf("⚠️ Could not retrieve user %s from Redis", userID)
		}
	}

	return users, nil
}

// UpdateUserChatID updates the user's chat ID in Redis
func (r *UserRepository) UpdateUserChatID(ctx context.Context, userID, chatID string) error {
	userKey := fmt.Sprintf("user:%s", userID)

	// ✅ Store ChatID in Redis
	_, err := r.client.HSet(ctx, userKey, "chatID", chatID).Result()
	if err != nil {
		log.Printf("❌ Error updating chat ID for user %s: %v", userID, err)
		return err
	}

	log.Printf("✅ Updated ChatID for user %s: %s", userID, chatID)
	return nil
}


// GetQueueLength returns the number of users in the waiting queue
func (r *UserRepository) GetQueueLength(ctx context.Context) (int, error) {
	count, err := r.client.ZCard(ctx, r.queue).Result()
	if err != nil {
		log.Printf("❌ Error getting queue length: %v", err)
		return 0, err
	}
	return int(count), nil
}
