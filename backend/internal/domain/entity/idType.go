// IDのタイプを分け、隠蔽するための構造体
// 現在はすべてのIDがUUIDを用いている。
// MySQLの内部でUUIDをBinaryで保存するために変換メソッドを用意している。
// 作成ロジックは/backend/internal/interface/factory/idFactory.goにinterfaceが定義され、
// /backend/internal/infrastructure/factoryImpl/idFactoryImpl.goでDI用の実装がされている。

package entity

import "github.com/google/uuid"

type UserID string
// UserID -> UUID変換メソッド
func (u *UserID) UserID2UUID() (uuid.UUID, error) {
	id, err := uuid.Parse(string(*u))
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
// UUID -> UserID変換メソッド
func (u *UserID) UUID2UserID(id uuid.UUID) {
	*u = UserID(id.String())
}

type RoomID string
// RoomID -> UUID変換メソッド
func (r *RoomID) RoomID2UUID() (uuid.UUID, error) {
	id, err := uuid.Parse(string(*r))
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
// UUID -> RoomID変換メソッド
func (r *RoomID) UUID2RoomID(id uuid.UUID) {
	*r = RoomID(id.String())
}

type MessageID string
// MessageID -> UUID変換メソッド
func (m *MessageID) MessageID2UUID() (uuid.UUID, error) {
	id, err := uuid.Parse(string(*m))
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
// UUID -> MessageID変換メソッド
func (m *MessageID) UUID2MessageID(id uuid.UUID) {
	*m = MessageID(id.String())
}

type WsClientID string
// WsClientID -> UUID変換メソッド
func (w *WsClientID) WsClientID2UUID() (uuid.UUID, error) {
	id, err := uuid.Parse(string(*w))
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
// UUID -> WsClientID変換メソッド
func (w *WsClientID) UUID2WsClientID(id uuid.UUID) {
	*w = WsClientID(id.String())
}
