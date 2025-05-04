package factoryimpl

import (
	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/interface/factory"
	"github.com/google/uuid"
)

type UserIDFactoryImpl struct{}

func NewUserIDFactory() factory.UserIDFactory {
	return &UserIDFactoryImpl{}
}

func (f *UserIDFactoryImpl) NewUserID() (entity.UserID, error) {
	return entity.UserID(uuid.New().String()), nil
}

type RoomIDFactoryImpl struct{}

func NewRoomIDFactory() factory.RoomIDFactory {
	return &RoomIDFactoryImpl{}
}

func (f *RoomIDFactoryImpl) NewRoomID() (entity.RoomID, error) {
	return entity.RoomID(uuid.New().String()), nil
}

type MessageIDFactoryImpl struct{}

func NewMessageIDFactory() factory.MessageIDFactory {
	return &MessageIDFactoryImpl{}
}

func (f *MessageIDFactoryImpl) NewMessageID() (entity.MessageID, error) {
	return entity.MessageID(uuid.New().String()), nil
}

type WsClientIDFactoryImpl struct{}

func NewWsClientIDFactory() factory.WsClientIDFactory {
	return &WsClientIDFactoryImpl{}
}

func (f *WsClientIDFactoryImpl) NewWsClientID() (entity.WsClientID, error) {
	return entity.WsClientID(uuid.New().String()), nil
}
