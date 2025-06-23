package gorillawebsocket

import (
	"errors"
	"time"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/domain/service"
	"example.com/infrahandson/internal/interface/adapter"
)

type GorillaWebSocketConnection struct {
	conn adapter.ConnAdapter
}

type NewGorillaWebSocketConnectionParams struct {
	Conn adapter.ConnAdapter
}

func (p *NewGorillaWebSocketConnectionParams) Validate() error {
	if p.Conn == nil {
		return errors.New("conn is required")
	}
	return nil
}

func NewGorillaWebSocketConnection(
	p *NewGorillaWebSocketConnectionParams,
) service.WebSocketConnection {
	if err := p.Validate(); err != nil {
		panic(err)
	}
	return &GorillaWebSocketConnection{
		conn: p.Conn,
	}
}

type MessageDTO struct {
	ID       entity.MessageID       // メッセージID
	RoomID   entity.RoomID          // 所属するチャットルームのID
	UserID   entity.UserID          // 投稿者のID（匿名なら名前など）
	Content  string                 // 本文
	SentAt   time.Time              // 送信日時
}

func (m *MessageDTO) ToEntity() *entity.Message {
	return entity.NewMessage(entity.MessageParams{
		ID:       m.ID,
		RoomID:   m.RoomID,
		UserID:   m.UserID,
		Content:  m.Content,
		SentAt:   m.SentAt,
	})
}

func (m *MessageDTO) FromEntity(msg *entity.Message) {
	m.ID = msg.GetID()
	m.ID = msg.GetID()
	m.RoomID = msg.GetRoomID()
	m.UserID = msg.GetUserID()
	m.Content = msg.GetContent()
	m.SentAt = msg.GetSentAt()
}

func (c *GorillaWebSocketConnection) ReadMessage() (*entity.Message, error) {
	var msgDTO MessageDTO
	err := c.conn.ReadJSON(&msgDTO)
	if err != nil {
		return nil, err
	}

	return msgDTO.ToEntity(), nil
}

func (c *GorillaWebSocketConnection) WriteMessage(msg *entity.Message) error {
	msgDTO := MessageDTO{}
	msgDTO.FromEntity(msg)
	err := c.conn.WriteJSON(msgDTO)
	if err != nil {
		return err
	}
	return nil
}

func (c *GorillaWebSocketConnection) Close() error {
	err := c.conn.CloseFunc()
	if err != nil {
		return err
	}
	return nil
}
