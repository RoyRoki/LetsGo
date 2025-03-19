package router

import (
	"github.com/gorilla/mux"
	"github.com/royroki/LetsGo/internal/modules/chat/presentation/controller"
)

func SetupChatRouter(chatController *controller.ChatController) *mux.Router {
	router := mux.NewRouter()

	// WebSocket route for chat
	router.HandleFunc("/ws", chatController.HandleConnection).Methods("GET")

	return router
}
