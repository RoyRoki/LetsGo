package controller

import (
	"net/http"

	"github.com/royroki/LetsGo/internal/modules/chat/application/interfaces"
	web_socket "github.com/royroki/LetsGo/internal/modules/chat/presentation/websocket"
)

type ChatController struct {
	chatUseCase      interfaces.ChatUseCase
	webSocketHandler *web_socket.WebSocketHandler
}

func NewChatController(chatUseCase interfaces.ChatUseCase, wsHandler *web_socket.WebSocketHandler) *ChatController {
	return &ChatController{
		chatUseCase:      chatUseCase,
		webSocketHandler: wsHandler,
	}
}

func (c *ChatController) HandleConnection(w http.ResponseWriter, r *http.Request) {
	c.webSocketHandler.HandleWSConnection(w, r)
}
