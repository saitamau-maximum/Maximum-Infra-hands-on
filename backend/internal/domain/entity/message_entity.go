package entity

import "time"

type Message struct {
	id      MessageID // メッセージID
	roomID  RoomID // 所属するチャットルームのID
	userID  UserID // 投稿者のID（匿名なら名前など）
	content string    // 本文
	sentAt  time.Time // 送信日時
}
