package domain

import "context"

type MatchRequest struct {
	UserID int32
	Tags   []string
	Module string // "chat", "call", "video"
}

type MatchResult struct {
	PartnerID int32
	Matched   bool
}

type Matcher interface {
	Match(ctx context.Context, req MatchRequest) (res MatchRequest, err error)
}
