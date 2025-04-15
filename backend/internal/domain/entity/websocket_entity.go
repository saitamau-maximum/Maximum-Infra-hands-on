package entity

import "time"

// factory関数はinterface層で定義されている
type Message struct {
	id      MessageID // メッセージID
	roomID  RoomID // 所属するチャットルームのID
	userID  UserID // 投稿者のID（匿名なら名前など）
	content string    // 本文
	sentAt  time.Time // 送信日時
}

// ID生成方法を隠蔽するためのラッパー
type WebsicketClientID string

// factory関数はinterface層で定義されている
type WebsocketClient struct {
	id       WebsicketClientID
	roomID   RoomID
	userName string
}


