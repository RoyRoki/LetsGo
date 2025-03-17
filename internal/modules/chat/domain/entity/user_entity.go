package entity

import "time"

// User represents a connected user in the chat system
type User struct {
	UserID   string    `json:"user_id"`
	ConnID   string    `json:"conn_id"`
	JoinTime time.Time `json:"join_time"`
	Chatted  int64     `json:"chatted"`
}
