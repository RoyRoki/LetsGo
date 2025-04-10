package entity

import enum "github.com/royroki/LetsGo/internal/common/matching/domain/enums"

type MatchRequest struct {
	UserID uint32
	Tags   []string
	Module enum.Module
}
