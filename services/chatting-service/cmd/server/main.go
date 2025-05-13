package main

import (
	"log"
	"net/http"
	"os"

	"github.com/royroki/services/chatting-service/internal/app/usecases"
	"github.com/royroki/services/chatting-service/internal/constants"
	grpcclient "github.com/royroki/services/chatting-service/internal/infrastructure/grpc"
	"github.com/royroki/services/chatting-service/internal/infrastructure/redis"
	"github.com/royroki/services/chatting-service/internal/infrastructure/websocket"
	"github.com/royroki/services/chatting-service/internal/presentation/controller"
	"github.com/royroki/services/chatting-service/internal/presentation/router"
)

func main() {
	// Initialize Redis
	redis.InitRedis()

	// Load gRPC matcher service address and port from environment
	grpcAddr := os.Getenv(constants.GRPCADD)
	grpcPort := os.Getenv(constants.GRPCPORT)

	if grpcAddr == "" || grpcPort == "" {
		log.Fatal("Environment variables MATCHER_GRPC_ADDR or MATCHER_GRPC_PORT are not set")
	}

	grpcTarget := grpcAddr + ":" + grpcPort
	log.Printf("Connecting to matcher service at %s", grpcTarget)

	// Initialize gRPC matcher client
	matcherClient := grpcclient.NewMatcherClient(grpcTarget)
	// Init WebSocket server
	wsServer := websocket.NewWCServer()

	// Init Usecase
	chatUsecase := usecases.NewChatUsecase(wsServer, matcherClient)

	// Init chat controller
	chatController := controller.NewChatController(wsServer, chatUsecase)

	// Setup router with dependencies
	r := router.SetupRoutes(chatController)

	// Start HTTP server
	port := os.Getenv(constants.WSPORT)
	if port == "" {
		port = "8080"
	}
	log.Printf("🚀 Server running on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
