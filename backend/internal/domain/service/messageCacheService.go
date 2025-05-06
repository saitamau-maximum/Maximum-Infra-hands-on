package service

import "example.com/infrahandson/internal/domain/entity"

const RECENT_MESSAGE_LIMIT = 20

type MessageCacheService interface {
	// 指定したルームの最近のメッセージを取得（RECENT_MESSAGE_LIMIT件）
	GetRecentMessages(roomID entity.RoomID) ([]*entity.Message, error)
	// メッセージをキャッシュに追加（RECENT_MESSAGE_LIMIT件を超えた場合は古いものから削除）
	AddMessage(roomID entity.RoomID, message *entity.Message) error
}

// デフォルトの最近のメッセージLimit数を取得
func DefaultRecentMessageLimit() int {
	return RECENT_MESSAGE_LIMIT
}
