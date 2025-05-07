package grpc

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/royroki/letsgo/services/matching-service/internal/app/usecases"
	"github.com/royroki/letsgo/services/matching-service/internal/constants"
	"github.com/royroki/letsgo/services/matching-service/internal/domain"
	pb "github.com/royroki/letsgo/services/matching-service/proto/matching"
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
	result, err := s.usecase.TryMatch(ctx, domain.MatchRequest{
		UserID: int32(req.UserId),
		Module: req.Module,
		Tags:   req.Tags,
	})
	log.Printf("Match User Called %s, %b", result.PartnerID, result.Matched)
	if err != nil {
		return nil, err
	}
	return &pb.MatchResponse{
		PartnerId: int64(result.PartnerID),
		Matched:   result.Matched,
	}, nil
}
