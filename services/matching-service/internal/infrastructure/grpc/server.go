package grpc

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/royroki/letsgo/services/matching-service/internal/app/usecases"
	"github.com/royroki/letsgo/services/matching-service/internal/constants"
	"github.com/royroki/letsgo/services/matching-service/internal/domain"
	"github.com/royroki/letsgo/services/matching-service/internal/infrastructure/redis"
	"github.com/royroki/letsgo/services/matching-service/proto/matching"
	"google.golang.org/grpc"
)

type gRPCServer struct {
	matching.UnimplementedMatcherServer
	usecase *usecases.MatchUsercase
}

func NewGRPCServer() {
	port := os.Getenv(constants.GRPCPORT)
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Printf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	matching.RegisterMatcherServer(s, &gRPCServer{usecase: usecases.NewMatchUsecase()})
	log.Printf("gRPC server (matching-service) listening on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (s *gRPCServer) MatchUser(ctx context.Context, req *matching.MatchRequest) (*matching.MatchResponse, error) {
	if len(req.Tags) <= 0 {
		partnerID, isMatched := redis.FifoMatch(req.Module, req.UserId)
		log.Printf("User %d matched via FIFO with partner %d", req.UserId, partnerID)

		return &matching.MatchResponse{
			PartnerId: int32(partnerID),
			Matched:   isMatched,
		}, nil

	}
	result, err := s.usecase.TryMatch(ctx, domain.MatchRequest{
		UserID: req.UserId,
		Module: req.Module,
		Tags:   req.Tags,
	})
	log.Printf("Match result for user %d: Matched=%v, PartnerID=%d", req.UserId, result.Matched, result.PartnerID)
	if err != nil {
		return nil, err
	}
	return &matching.MatchResponse{
		PartnerId: result.PartnerID,
		Matched:   result.Matched,
	}, nil
}
