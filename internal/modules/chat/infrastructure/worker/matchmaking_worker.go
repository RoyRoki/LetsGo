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
// Run starts the matchmaking loop
func (w *MatchmakingWorker) Run() {
	log.Println("🔄 Matchmaking Worker Started...")

	for {
		ctx := context.Background()

		// ✅ Step 1: Pop users (but only if 2+ exist)
		users, err := w.userRepo.PopTopUsers(ctx, 2)
		if err != nil {
			log.Printf("❌ Error retrieving users from queue: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// ✅ Step 2: If only 1 user remains, re-add them
		if len(users) == 1 {
			log.Printf("⚠️ Only one user (%s) left in queue. Re-adding them...", users[0].UserID)
			if err := w.userRepo.AddUserToQueue(ctx, users[0]); err != nil {
				log.Printf("❌ Failed to re-add user to queue: %v", err)
			}
			time.Sleep(5 * time.Second)
			continue
		}

		// ✅ Step 3: Pair users (always in groups of 2)
		for len(users) >= 2 {
			userA, userB := users[0], users[1]
			users = users[2:]

			if err := w.chatUsecase.HandleChatPair(ctx, userA, userB); err != nil {
				log.Printf("❌ Failed to pair users %s & %s: %v", userA.UserID, userB.UserID, err)
				continue
			}

			log.Printf("✅ Matched Users: %s <-> %s", userA.UserID, userB.UserID)
			w.chatUsecase.ListenFromConnection(userA.ConnID)
			w.chatUsecase.ListenFromConnection(userB.ConnID)
		}

		// ✅ Step 4: Sleep before next matchmaking check
		time.Sleep(5 * time.Second)
	}
}
