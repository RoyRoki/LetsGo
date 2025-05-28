package usecases

import (
	"log"

	"github.com/royroki/services/chatting-service/internal/constants"
	"github.com/royroki/services/chatting-service/internal/domain/entity"
	"github.com/royroki/services/chatting-service/internal/domain/repo"
	grpcclient "github.com/royroki/services/chatting-service/internal/infrastructure/grpc"
)

type ChatUsecase struct {
	wsServer      repo.ChatGateway
	matcherClient *grpcclient.MatcherClient
}

func NewChatUsecase(ws repo.ChatGateway, matcher *grpcclient.MatcherClient) *ChatUsecase {
	return &ChatUsecase{
		wsServer:      ws,
		matcherClient: matcher,
	}
}

// HandleMatchRequest processes user's tags and attempts to find a match
func (cu *ChatUsecase) HandleMatchRequest(user *entity.User) {
tryMatch:
	partnerID, matched, err := cu.matcherClient.Match(user.ID, user.Module, user.Tags)
	if err != nil {
		log.Printf("gRPC matching error for user %d: %v", user.ID, err)
		cu.wsServer.Send(user.ID, entity.Message{
			From:    0,
			Content: "Server error during matching",
			Status:  constants.ServerError,
		})
		return
	}

	if !matched {
		cu.wsServer.Send(user.ID, entity.Message{
			From:    0,
			Content: "No match found. Waiting...",
			Status:  constants.MatchNotFound,
		})
		return
	}

	// Check if partner is active
	if !cu.wsServer.IsActive(partnerID) {
		log.Printf("Partner %d is inactive, re-matching...", partnerID)
		goto tryMatch
	}

	// Check if partner is already paired with someone else
	if cu.wsServer.GetPartnerID(partnerID) != 0 {
		log.Printf("Partner %d is already paired, re-matching...", partnerID)
		goto tryMatch
	}

	// Check if user itself is already paired (optional safety check)
	if cu.wsServer.GetPartnerID(user.ID) != 0 {
		log.Printf("User %d is already paired, skipping match request", user.ID)
		return
	}

	// Pair users
	cu.wsServer.Pair(user.ID, partnerID)

	log.Printf("Matched users: %d <--> %d with tags: %v\n", user.ID, partnerID, user.Tags)

	// Notify both users
	cu.wsServer.Send(user.ID, entity.Message{
		From:    0,
		Content: "Connected to a partner!",
		Status:  constants.MatchFound,
	})

	cu.wsServer.Send(partnerID, entity.Message{
		From:    0,
		Content: "You have been connected!",
		Status:  constants.MatchFound,
	})
}


// HandleDisconnect cleans up and notifies partner
func (cu *ChatUsecase) HandleDisconnect(userID int32) {
	partnerID := cu.wsServer.GetPartnerID(userID)

	cu.wsServer.Unregister(userID)

	if partnerID != 0 {
		cu.wsServer.Send(partnerID, entity.Message{
			From:    0,
			Content: "Your partner has disconnected.",
			Status:  constants.PartnerDisconnected,
		})

		// Optionally requeue the partner
		if partner := cu.wsServer.GetUser(partnerID); partner != nil {
			go cu.HandleMatchRequest(partner)
		}
	}
}
