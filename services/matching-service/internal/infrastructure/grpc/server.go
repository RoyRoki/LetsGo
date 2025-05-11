package grpc

import (
	"context"
	"log"
	"net"
	"os"

	pb "github.com/royroki/letsgo/proto/matching"
	"github.com/royroki/letsgo/services/matching-service/internal/app/usecases"
	"github.com/royroki/letsgo/services/matching-service/internal/constants"
	"github.com/royroki/letsgo/services/matching-service/internal/domain"
	"github.com/royroki/letsgo/services/matching-service/internal/infrastructure/redis"
	"google.golang.org/grpc"
)

type gRPCServer struct {
	pb.UnimplementedMatcherServer
	usecase *usecases.MatchUsercase
}

func NewGRPCServer() {
	port := os.Getenv(constants.GRPCPORT)
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Printf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterMatcherServer(s, &gRPCServer{usecase: usecases.NewMatchUsecase()})
	log.Printf("gRPC server (matching-service) listening on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (s *gRPCServer) MatchUser(ctx context.Context, req *pb.MatchRequest) (*pb.MatchResponse, error) {
	if len(req.Tags) <= 0 {
		partnerID, isMatched := redis.FifoMatch(req.Module, int32(req.UserId))
		log.Printf("User %d matched via FIFO with partner %d", req.UserId, partnerID)
		return &pb.MatchResponse{
			PartnerId: int64(partnerID),
			Matched:   isMatched,
		}, nil

	}
	result, err := s.usecase.TryMatch(ctx, domain.MatchRequest{
		UserID: int32(req.UserId),
		Module: req.Module,
		Tags:   req.Tags,
	})
	log.Printf("Match result for user %d: Matched=%v, PartnerID=%d", req.UserId, result.Matched, result.PartnerID)
	if err != nil {
		return nil, err
	}
	return &pb.MatchResponse{
		PartnerId: int64(result.PartnerID),
		Matched:   result.Matched,
	}, nil
}
