package worker

import (
	"context"
	"log"
	"time"

	"github.com/royroki/LetsGo/internal/modules/chat/application/usecase"
	"github.com/royroki/LetsGo/internal/modules/chat/domain/repository"
)

// MatchmakingWorker continuously pairs users from the queue
func MatchmakingWorker(chatUsecase *usecase.ChatUseCase, userRepo repository.UserRepository) {
	log.Println("🔄 Matchmaking Worker Started...")

	for {
		ctx := context.Background()

		// Pop the top two users from the queue
		users, err := userRepo.PopTopUsers(ctx, 2) // Atomic removal
		if err != nil || len(users) < 2 {
			log.Println("⚠️ Not enough users to match. Waiting...")
			time.Sleep(5 * time.Second)
			continue
		}

		// Pair users (Since we pop two at a time, always pair them)
		for len(users) >= 2 {
			userA, userB := users[0], users[1] // Get two users
			users = users[2:]                  // Remove paired users from the list

			// Pair and create a chat session
			err := chatUsecase.HandleChatPair(ctx, userA, userB)
			if err != nil {
				log.Println("⚠️ Failed to pair users:", err)
				continue
			}

			log.Printf("✅ Matched Users: %s <-> %s", userA.UserID, userB.UserID)
		}

		// Step 3️⃣: Sleep before checking again
		time.Sleep(5 * time.Second)
	}
}
