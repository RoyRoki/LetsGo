package controller

import (
	"net/http"

	"github.com/royroki/LetsGo/internal/modules/chat/application/usecase"
	"github.com/royroki/LetsGo/internal/modules/chat/presentation/websocket"
)

type ChatController struct {
	chatUseCase *usecase.ChatUseCase
	connections map[string]*websocket.WebSocketHub
}

func NewChatController(chatUseCase *usecase.ChatUseCase) *ChatController {
	return &ChatController{chatUseCase: chatUseCase}
}

func (c *ChatController) HandleConnection(w http.ResponseWriter, r *http.Request) {

}
