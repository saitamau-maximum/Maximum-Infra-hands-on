// IDの生成方法を隠蔽するためのラッパー
package entity

import "github.com/google/uuid"

type UserID string

func (u *UserID) UserID2UUID() (uuid.UUID, error) {
	id, err := uuid.Parse(string(*u))
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (u *UserID) UUID2UserID(id uuid.UUID) {
	*u = UserID(id.String())
}

type RoomID string

func (r *RoomID) RoomID2UUID() (uuid.UUID, error) {
	id, err := uuid.Parse(string(*r))
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
func (r *RoomID) UUID2RoomID(id uuid.UUID) {
	*r = RoomID(id.String())
}

type MessageID string

func (m *MessageID) MessageID2UUID() (uuid.UUID, error) {
	id, err := uuid.Parse(string(*m))
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
func (m *MessageID) UUID2MessageID(id uuid.UUID) {
	*m = MessageID(id.String())
}

type WsClientID string

func (w *WsClientID) WsClientID2UUID() (uuid.UUID, error) {
	id, err := uuid.Parse(string(*w))
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
func (w *WsClientID) UUID2WsClientID(id uuid.UUID) {
	*w = WsClientID(id.String())
}