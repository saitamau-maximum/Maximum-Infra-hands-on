package entity

import "time"

type Message struct {
	id      MessageID // メッセージID
	roomID  RoomID    // 所属するチャットルームのID
	userID  UserID    // 投稿者のID（匿名なら名前など）
	content string    // 本文
	sentAt  time.Time // 送信日時
}

type MessageParams struct {
	ID      MessageID
	RoomID  RoomID
	UserID  UserID
	Content string
	SentAt  time.Time
}

func NewMessage(params MessageParams) *Message {
	return &Message{
		id:      params.ID,
		roomID:  params.RoomID,
		userID:  params.UserID,
		content: params.Content,
		sentAt:  params.SentAt,
	}
}

func (m *Message) GetID() MessageID {
	return m.id
}

func (m *Message) GetRoomID() RoomID {
	return m.roomID
}
func (m *Message) GetUserID() UserID {
	return m.userID
}

func (m *Message) GetContent() string {
	return m.content
}

func (m *Message) GetSentAt() time.Time {
	return m.sentAt
}
