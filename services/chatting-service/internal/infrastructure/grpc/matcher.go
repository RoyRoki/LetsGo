package grpcclient

import (
	"context"
	"log"
	"time"

	"github.com/royroki/letsgo/proto/matching"
	"google.golang.org/grpc"
)

type MatcherClient struct {
	client matching.MatcherClient
}

func NewMatcherClient(address string) *MatcherClient {
	conn, err := grpc.Dial(address, grpc.WithInsecure()) // Use TLS in prod
	if err != nil {
		log.Fatalf("Failed to connect to matcher service: %v", err)
	}

	client := matching.NewMatcherClient(conn)
	return &MatcherClient{client: client}
}

func (mc *MatcherClient) Match(userID int64, module string, tags []string) (int64, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := mc.client.MatchUser(ctx, &matching.MatchRequest{
		UserId: userID,
		Module: module,
		Tags:   tags,
	})
	if err != nil {
		return 0, false, err
	}

	return resp.PartnerId, resp.Matched, nil
}
