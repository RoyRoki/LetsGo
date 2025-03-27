package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	// Use interface instead of concrete implementation
	var chatUsecase interfaces.ChatUseCase = usecase.NewChatUseCase(chatService)

	wsHandler := web_socket.NewWebSocketHandler(chatUsecase, wsHub)

	chatController := controller.NewChatController(chatUsecase, wsHandler)

	chatRouter := router.SetupChatRouter(chatController)

	// Start worker
	chatWorker := worker.NewMatchmakingWorker(chatUsecase, userRepo)

	go chatWorker.Run()

	server := &http.Server{
		Addr:    ":8080",
		Handler: chatRouter,
	}

	// Start WebSocket Server
	go func() {
		log.Println("✅ WebSocket Server started at ws://localhost:8080/ws")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server error: %v", err)
		}
	}()

	// Wait for termination signal (Ctrl+C, SIGTERM)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	log.Println("🚀 Shutting down server...")

	// Cleanup resources
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Stop background worker
	chatWorker.Stop() 
	// Gracefully shutdown the HTTP server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("❌ HTTP Server Shutdown Failed: %v", err)
	}

	// Stop WebSocket connections
	wsHub.Shutdown() // Implement a `Shutdown` method in `WebSocketHub` to clean connections.

	if err = redisClient.Del(ctx, "waiting_queue").Err(); err != nil {
		log.Println("Failed to clear waiting_queue")
	}
	log.Println("waiting_queue deleted")

	// Close Redis connection
	if err := redisClient.Close(); err != nil {
		log.Fatalf("❌ Redis close error: %v", err)
	}

	log.Println("✅ Server shutdown complete")

}
