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
}

// NewMatchmakingWorker initializes a MatchmakingWorker
func NewMatchmakingWorker(chatUsecase interfaces.ChatUseCase, userRepo repository.UserRepository) *MatchmakingWorker {
	return &MatchmakingWorker{
		chatUsecase: chatUsecase,
		userRepo:    userRepo,
	}
}

// Run starts the matchmaking loop
func (w *MatchmakingWorker) Run() {
	log.Println("🔄 Matchmaking Worker Started...")

	for {
		ctx := context.Background()

		// Pop the top two users from the queue
		users, err := w.userRepo.PopTopUsers(ctx, 2) // Atomic removal
		if err != nil || len(users) < 2 {
			time.Sleep(5 * time.Second)
			log.Printf("Total waiting user : %d", len(users))
			continue
		}

		// Pair users
		for len(users) >= 2 {
			userA, userB := users[0], users[1]
			users = users[2:] // Remove paired users from the list

			// Handle chat pairing
			if err := w.chatUsecase.HandleChatPair(ctx, userA, userB); err != nil {
				log.Println("⚠️ Failed to pair users:", err)
				continue
			}

			log.Printf("✅ Matched Users: %s <-> %s", userA.UserID, userB.UserID)
		}

		// Sleep before the next check
		time.Sleep(5 * time.Second)
	}
}
