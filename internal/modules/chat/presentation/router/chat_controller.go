package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/royroki/LetsGo/internal/modules/chat/infrastructure/waitingqueue"
	"github.com/royroki/LetsGo/internal/modules/chat/presentation/websocket"
)

func SetupRouter(wq *waitingqueue.WaitingQueue) *mux.Router {
	router := mux.NewRouter()
	hub := websocket.NewWebSocketHub(wq)

	// WebSocket route for chat
	router.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		hub.HandleWebSocket(w, r)
	}).Methods("GET")

	return router
}
