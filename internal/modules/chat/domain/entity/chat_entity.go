package entity

import "time"

// Chat represents a conversation session between two users
type Chat struct {
	ID        string     `json:"id"`
	UserA     User       `json:"user_a"`
	UserB     User       `json:"user_b"`
	StartTime time.Time  `json:"start_time"`
	EndTime   *time.Time `json:"end_time,omitempty"` // Pointer to handle ongoing chats (nil if active)
}
