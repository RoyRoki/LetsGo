package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/royroki/LetsGo/internal/common/logger"
	"github.com/royroki/LetsGo/internal/config/database"
	config "github.com/royroki/LetsGo/internal/config/env"
	"github.com/royroki/LetsGo/internal/modules/chat/infrastructure/waitingqueue"
	"github.com/royroki/LetsGo/internal/modules/chat/presentation/router"
)

func main() {
	config.NewEnvConfig()
	logger.NewZapLogger()

	redis_config := database.NewRedisConfig()
	err := redis_config.Ping()
	if err != nil {
		print("Failed To Ping Redis")
	}

	rdb := redis_config.NewClient()

	queue := waitingqueue.NewWaitingQueue(rdb, "chatQueue")

	r := router.SetupRouter(queue)

	// Start the HTTP server
	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", r))
}
