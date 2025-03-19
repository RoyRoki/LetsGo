package main

import (
	"log"
	"net/http"

	"github.com/royroki/LetsGo/internal/common/logger"
	"github.com/royroki/LetsGo/internal/config/database"
	config "github.com/royroki/LetsGo/internal/config/env"
	"github.com/royroki/LetsGo/internal/modules/chat/application/interfaces"
	"github.com/royroki/LetsGo/internal/modules/chat/application/usecase"
	service "github.com/royroki/LetsGo/internal/modules/chat/domain/services"
	"github.com/royroki/LetsGo/internal/modules/chat/infrastructure/persistence"
	web_socket_hub "github.com/royroki/LetsGo/internal/modules/chat/infrastructure/websocket"
	"github.com/royroki/LetsGo/internal/modules/chat/infrastructure/worker"
	"github.com/royroki/LetsGo/internal/modules/chat/presentation/controller"
	"github.com/royroki/LetsGo/internal/modules/chat/presentation/router"
	web_socket "github.com/royroki/LetsGo/internal/modules/chat/presentation/websocket"
)

func main() {
	config.NewEnvConfig()
	logger.NewZapLogger()

	redisConfig := database.NewRedisConfig()
	err := redisConfig.Ping()
	if err != nil {
		print("Failed To Ping Redis")
	}

	redisClient := redisConfig.NewClient()

	// r := router.SetupRouter(queue)

	wsHub := web_socket_hub.NewWebSocketHub()
	userRepo := persistence.NewUserRepository(redisClient, "waiting_queue")
	chatRepo := persistence.NewChatRepository(redisClient)

	chatService := service.NewChatService(chatRepo, userRepo, wsHub)

	// Start worker
	go worker.MatchmakingWorker()

	// Use interface instead of concrete implementation
	var chatUsecase interfaces.ChatUseCase = usecase.NewChatUseCase(chatService)

	wsHandler := web_socket.NewWebSocketHandler(chatUsecase, wsHub)

	chatController := controller.NewChatController(chatUsecase, wsHandler)

	chatRouter := router.SetupChatRouter(chatController)

	// Start the HTTP server
	log.Println("✅ WebSocket Server started at ws://localhost:8080/ws")
	log.Fatal(http.ListenAndServe(":8080", chatRouter))

}
