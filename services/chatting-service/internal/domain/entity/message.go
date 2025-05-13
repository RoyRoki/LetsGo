package entity

import (
	"encoding/json"
)

type Message struct {
	From      int64  `json:"from"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
	Status    string `json:"status"` // status according the source
}

// ToJSON converts the Message to a JSON string
func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}
