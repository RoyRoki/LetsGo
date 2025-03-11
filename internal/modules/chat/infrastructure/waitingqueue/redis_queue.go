// internal/modules/chat/infrastructure/redis_queue.go
package waitingqueue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"slices"

	"github.com/go-redis/redis/v8"
)

type UserAttributes struct {
	Tags   []string `json:"tags"`
	Gender string   `json:"gender"`
}

type WaitingQueue struct {
	client *redis.Client
	queue  string
}

func NewWaitingQueue(client *redis.Client, queueName string) *WaitingQueue {
	return &WaitingQueue{
		client: client,
		queue:  queueName,
	}
}

// AddUserToQueue adds a user to the waiting queue.
func (wq *WaitingQueue) AddUserToQueue(userID string) error {
	ctx := context.Background()
	// Add the user to the queue in Redis
	_, err := wq.client.LPush(ctx, wq.queue, userID).Result()
	if err != nil {
		log.Printf("Error adding user to queue: %v", err)
		return err
	}
	log.Printf("User %s added to the queue", userID)
	return nil
}

// GetNextUser gets the next available user from the queue.
func (wq *WaitingQueue) GetNextUser(userID string) (string, error) {
	ctx := context.Background()

	// Attempt to retrieve the next user from the queue
	user, err := wq.client.RPop(ctx, wq.queue).Result()
	if err != nil {
		log.Printf("Error retrieving user from queue: %v", err)
		return "", err
	}

	// If the retrieved user is the same as the current user, put them back and continue
	if user == userID {
		log.Printf("Retrieved the same user (%s), re-adding to the queue and trying again.", userID)
		// Add the user back into the queue to continue waiting
		if err := wq.AddUserToQueue(userID); err != nil {
			log.Printf("Error adding user %s back to queue: %v", userID, err)
			return "", err
		}
		// Recursively call GetNextUser to get another user
		return wq.GetNextUser(userID)
	}

	log.Printf("Next user retrieved from queue: %s", user)
	return user, nil
}

// // PairUsers checks if there are at least two users in the queue, pairs them, and returns their IDs.
// func (wq *WaitingQueue) PairUsers() (string, string, error) {
// 	context.Background()

// 	// Ensure at least two users are in the queue before pairing
// 	user1, err := wq.GetNextUser()
// 	if err != nil {
// 		return "", "", err
// 	}

// 	user2, err := wq.GetNextUser()
// 	if err != nil {
// 		// If only one user is available, put them back in the queue
// 		// You can also add a waiting message to the user
// 		_ = wq.AddUserToQueue(user1)
// 		log.Printf("Only one user found, putting %s back in the queue", user1)
// 		return "", "", nil
// 	}

// 	log.Printf("Users %s and %s have been paired", user1, user2)
// 	return user1, user2, nil
// }

// AddUserWithAttributes adds user attributes to Redis
func (wq *WaitingQueue) AddUserWithAttributes(userID string, attrs UserAttributes) error {
	attrsJSON, _ := json.Marshal(attrs)
	return wq.client.HSet(context.Background(), "user:"+userID, "attributes", attrsJSON).Err()
}

// GetQueueLength returns the length of the waiting queue
func (wq *WaitingQueue) GetQueueLength() (int64, error) {
	return wq.client.LLen(context.Background(), wq.queue).Result()
}

// GetNextUserWithCriteria retrieves a user that matches some pairing criteria
func (wq *WaitingQueue) GetNextUserWithCriteria(criteria UserAttributes) (string, error) {
	users, err := wq.client.LRange(context.Background(), wq.queue, 0, -1).Result()
	if err != nil {
		return "", err
	}

	for _, userID := range users {
		attrsJSON, err := wq.client.HGet(context.Background(), "user:"+userID, "attributes").Result()
		if err != nil {
			continue
		}

		var userAttrs UserAttributes
		err = json.Unmarshal([]byte(attrsJSON), &userAttrs)
		if err != nil {
			continue
		}

		if containsMatchingTags(userAttrs, criteria) {
			return userID, nil
		}
	}

	return "", fmt.Errorf("no matching user found")
}

func containsMatchingTags(userAttrs, criteria UserAttributes) bool {
	for _, tag := range criteria.Tags {
		if slices.Contains(userAttrs.Tags, tag) {
			return true
		}
	}
	return false
}
