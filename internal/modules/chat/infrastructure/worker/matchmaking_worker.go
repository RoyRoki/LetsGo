package worker

import (
	"context"
	"log"
	"time"

	"github.com/royroki/LetsGo/internal/modules/chat/application/interfaces"
	"github.com/royroki/LetsGo/internal/modules/chat/domain/repository"
)

// MatchmakingWorker handles user pairing from the queue
type MatchmakingWorker struct {
	chatUsecase interfaces.ChatUseCase
	userRepo    repository.UserRepository
	stopChan    chan struct{} // Stop signal channel
}

// NewMatchmakingWorker initializes a MatchmakingWorker
func NewMatchmakingWorker(chatUsecase interfaces.ChatUseCase, userRepo repository.UserRepository) *MatchmakingWorker {
	return &MatchmakingWorker{
		chatUsecase: chatUsecase,
		userRepo:    userRepo,
		stopChan:    make(chan struct{}),
	}
}

// Run starts the matchmaking loop
func (w *MatchmakingWorker) Run() {
	log.Println("🔄 Matchmaking Worker Started...")

	for {
		select {
		case <-w.stopChan:
			log.Println("🛑 Matchmaking Worker Stopped.")
			break
			
		default:
			ctx := context.Background()

			// ✅ Step 1: Check if at least 2 users exist before popping
			userCount, err := w.userRepo.GetQueueLength(ctx)
			if err != nil {
				log.Printf("❌ Error checking queue length: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}

			if userCount < 2 {
				log.Println("⚠️ Not enough users in queue, waiting...")
				time.Sleep(5 * time.Second)
				continue
			}

			// ✅ Step 2: Pop exactly 2 users
			users, err := w.userRepo.PopTopUsers(ctx, 2)
			if err != nil || len(users) != 2 {
				log.Printf("❌ Error retrieving users from queue: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}

			// ✅ Step 3: Pair users
			userA, userB := users[0], users[1]

			if err := w.chatUsecase.HandleChatPair(ctx, userA, userB); err != nil {
				log.Printf("❌ Failed to pair users %s & %s: %v", userA.UserID, userB.UserID, err)
				continue
			}

			log.Printf("✅ Matched Users: %s <-> %s", userA.UserID, userB.UserID)
			w.chatUsecase.ListenFromConnection(userA.UserID)
			w.chatUsecase.ListenFromConnection(userB.UserID)

			// ✅ Step 4: Sleep before next matchmaking check
			time.Sleep(5 * time.Second)
		}
	}
}
// Stop signals the matchmaking worker to terminate
func (w *MatchmakingWorker) Stop() {
	log.Println("🚀 Stopping Matchmaking Worker...")
	close(w.stopChan) // Sends a stop signal
}

