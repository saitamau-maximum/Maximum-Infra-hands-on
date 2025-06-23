// メッセージキャッシュのロジックインターフェース
// 具体実装は/infrastructure/messageCacheImpl
package service

import (
	"context"

	"example.com/infrahandson/internal/domain/entity"
)


// RECENT_MESSAGE_LIMIT defines the maximum number of recent messages to be cached and retrieved for a room.
// This limit should align with any client-side constraints to ensure consistent message handling between
// the frontend and backend.
const RECENT_MESSAGE_LIMIT = 20

type MessageCacheService interface {
	// GetRecentMessage は指定したルームの最近のメッセージを取得（RECENT_MESSAGE_LIMIT件）
	GetRecentMessages(ctx context.Context, roomID entity.RoomID) ([]*entity.Message, error)

	// AddMessage はメッセージをキャッシュに追加（RECENT_MESSAGE_LIMIT件を超えた場合は古いものから削除）
	AddMessage(ctx context.Context, roomID entity.RoomID, message *entity.Message) error
}

// DefaultRecentMessageLimit はデフォルトの最近のメッセージLimit数を取得
func DefaultRecentMessageLimit() int {
	return RECENT_MESSAGE_LIMIT
}
