package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/royroki/LetsGo/internal/modules/chat/application/usecase"
	"github.com/royroki/LetsGo/internal/modules/chat/presentation/controller"
)

func SetupRouter(chatUsecase *usecase.ChatUseCase) *mux.Router {
	router := mux.NewRouter()

	// WebSocket route for chat
	router.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		controller.NewChatController(chatUsecase)
	}).Methods("GET")

	return router
}
