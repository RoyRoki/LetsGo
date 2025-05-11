package entity

type Message struct {
	From      int64
	Content   string
	Timestamp int64
	System    bool // true if from server
}
