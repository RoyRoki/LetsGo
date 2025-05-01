package domain

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
	Match(req MatchRequest) (res MatchRequest, err error)
}
