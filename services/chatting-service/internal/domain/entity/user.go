package entity

import "github.com/gorilla/websocket"

type User struct {
	ID     int32
	Conn   *websocket.Conn
	IP     string
	Tags   []string
	Module string
}
