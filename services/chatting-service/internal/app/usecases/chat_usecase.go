package usecases

import (
	"log"

	"github.com/royroki/services/chatting-service/internal/constants"
	"github.com/royroki/services/chatting-service/internal/domain/entity"
	grpcclient "github.com/royroki/services/chatting-service/internal/infrastructure/grpc"
	"github.com/royroki/services/chatting-service/internal/infrastructure/websocket"
)

type ChatUsecase struct {
	wsServer      *websocket.WCServer
	matcherClient *grpcclient.MatcherClient
}

func NewChatUsecase(ws *websocket.WCServer, matcher *grpcclient.MatcherClient) *ChatUsecase {
	return &ChatUsecase{
		wsServer:      ws,
		matcherClient: matcher,
	}
}

// HandleMatchRequest processes user's tags and attempts to find a match
func (cu *ChatUsecase) HandleMatchRequest(user *entity.User) {
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

	// Check if partner is still active
	if !cu.wsServer.IsActive(partnerID) {
		log.Printf("Partner %d is inactive, re-matching...", partnerID)
		// Optionally call Match again or notify user
		cu.HandleMatchRequest(user)
		return
	}

	// Send success message to both
	cu.wsServer.Pair(user.ID, partnerID)

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
func (cu *ChatUsecase) HandleDisconnect(userID int64) {
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
