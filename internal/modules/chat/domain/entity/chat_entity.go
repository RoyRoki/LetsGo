package entity

type Chat struct {
	ID           string // Unique identifier for the chat session
	Participant1 string // User 1 in the chat
	Participant2 string // User 2 in the chat
	Status       string // Active, Ended, etc.
	StartTime    int64  // Timestamp when the chat started
	EndTime      int64  // Timestamp when the chat ended
}
