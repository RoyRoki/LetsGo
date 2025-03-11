package controller

import (
	"encoding/json"
	"net/http"

	"github.com/royroki/LetsGo/internal/modules/chat/application/usecase"
)

type ChatController struct {
	chatUseCase *usecase.ChatUseCase
}

func NewChatController(chatUseCase *usecase.ChatUseCase) *ChatController {
	return &ChatController{chatUseCase: chatUseCase}
}

func (c *ChatController) CreateChat(w http.ResponseWriter, r *http.Request) {
	var participants struct {
		Participant1 string `json:"participant1"`
		Participant2 string `json:"participant2"`
	}

	if err := json.NewDecoder(r.Body).Decode(&participants); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	chat, err := c.chatUseCase.CreateChat(participants.Participant1, participants.Participant2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(chat)
}
