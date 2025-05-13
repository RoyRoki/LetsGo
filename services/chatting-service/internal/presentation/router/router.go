package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/royroki/services/chatting-service/internal/presentation/controller"
	"github.com/royroki/services/chatting-service/internal/presentation/middleware"
)

func SetupRoutes(controller *controller.ChatController) http.Handler {
	r := mux.NewRouter()

	// Register WebSocket handler
	r.HandleFunc("/ws/chat", controller.HandleWebSocket)

	// Apply CORS middleware
	return middleware.CORS()(r)
}
