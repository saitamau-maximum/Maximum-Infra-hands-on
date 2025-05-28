// チャットメッセージの永続化のメソッドを先に定義
package repository

import (
	"context"
	"time"

	"example.com/infrahandson/internal/domain/entity"
)

type MessageRepository interface {
	// CreateMessage は指定された Message を永続化します。
	// コンテキストがキャンセルされた場合や保存に失敗した場合はエラーを返します。
	CreateMessage(ctx context.Context, msg *entity.Message) error

	// GetMessageHistoryInRoom は指定された部屋IDのメッセージ履歴を、指定された時刻より前のものから取得します。
	// 結果にはメッセージ配列、次ページ取得用の時刻、次ページが存在するかのフラグ、エラーを含みます。
	GetMessageHistoryInRoom(ctx context.Context, roomID entity.RoomID, limit int, beforeSentAt time.Time) (messages []*entity.Message, nextBeforeSentAt time.Time, hasNext bool, err error)
}

