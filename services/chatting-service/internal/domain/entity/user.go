package entity

type User struct {
	ID      int32
	Conn    WebSocketConn
	IP      string
	Tags    []string
	Module  string
	partner *User // nil if not matched
}

type WebSocketConn interface {
	WriteJSON(v interface{}) error
	ReadJSON(v interface{}) error
	Close() error
}
