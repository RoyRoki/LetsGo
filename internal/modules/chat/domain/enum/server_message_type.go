package enum

type ServerMessageType string

const (
	ServerMessageConnected    ServerMessageType = "connected"
	ServerMessageDisconnected ServerMessageType = "disconnected"
	ServerMessageTyping       ServerMessageType = "typing"
	ServerMessageWaiting      ServerMessageType = "waiting"
	ServerMessageError        ServerMessageType = "error"
)
