
package repo

import "github.com/royroki/services/chatting-service/internal/domain/entity"

type ChatGateway interface {
	Send(userID int32, msg entity.Message)
	IsActive(userID int32) bool
	GetPartnerID(userID int32) int32
	GetUser(userID int32) *entity.User
	Pair(userA, userB int32)
	Unregister(userID int32)
}