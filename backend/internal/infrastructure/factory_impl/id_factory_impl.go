package factory_impl

import (
	"errors"

	"example.com/webrtc-practice/internal/domain/entity"
	"example.com/webrtc-practice/internal/interface/factory"
	"github.com/google/uuid"
)

type UserIDFactoryImpl struct{}

func NewUserIDFactory() factory.UserIDFactory {
	return &UserIDFactoryImpl{}
}

func (f *UserIDFactoryImpl) NewUserPublicID() (entity.UserPublicID, error) {
	return entity.UserPublicID(uuid.New().String()), nil
}

func (f *UserIDFactoryImpl) FromInt(i int) (entity.UserID, error) {
	if i <= 0 {
		return 0, errors.New("invalid integer for UserID")
	}
	return entity.UserID(i), nil
}

func (f *UserIDFactoryImpl) FromString(s string) (entity.UserPublicID, error) {
	if _, err := uuid.Parse(s); err != nil {
		return "", err
	}
	return entity.UserPublicID(s), nil
}

type RoomIDFactoryImpl struct{}

func NewRoomIDFactory() factory.RoomIDFactory {
	return &RoomIDFactoryImpl{}
}

func (f *RoomIDFactoryImpl) NewRoomPublicID() (entity.RoomPublicID, error) {
	return entity.RoomPublicID(uuid.New().String()), nil
}

func (f *RoomIDFactoryImpl) FromInt(i int) (entity.RoomID, error) {
	if i <= 0 {
		return 0, errors.New("invalid integer for RoomID")
	}
	return entity.RoomID(i), nil
}

func (f *RoomIDFactoryImpl) FromString(s string) (entity.RoomPublicID, error) {
	if _, err := uuid.Parse(s); err != nil {
		return "", err
	}
	return entity.RoomPublicID(s), nil
}

type MessageIDFactoryImpl struct{}

func NewMessageIDFactory() factory.MessageIDFactory {
	return &MessageIDFactoryImpl{}
}

func (f *MessageIDFactoryImpl) NewMessagePublicID() (entity.MessagePublicID, error) {
	return entity.MessagePublicID(uuid.New().String()), nil
}

func (f *MessageIDFactoryImpl) FromString(s string) (entity.MessagePublicID, error) {
	if _, err := uuid.Parse(s); err != nil {
		return "", err
	}
	return entity.MessagePublicID(s), nil
}

func (f *MessageIDFactoryImpl) FromInt(i int) (entity.MessagePublicID, error) {
	if i <= 0 {
		return "", errors.New("invalid integer for MessageID")
	}
	return entity.MessagePublicID(uuid.New().String()), nil
}

type WsClientIDFactoryImpl struct{}

func NewWsClientIDFactory() factory.WsClientIDFactory {
	return &WsClientIDFactoryImpl{}
}

func (f *WsClientIDFactoryImpl) NewWsClientPublicID() (entity.WsClientPublicID, error) {
	return entity.WsClientPublicID(uuid.New().String()), nil
}

func (f *WsClientIDFactoryImpl) FromInt(i int) (entity.WsClientPublicID, error) {
	if i <= 0 {
		return "", errors.New("invalid integer for WsClientID")
	}
	return entity.WsClientPublicID(uuid.New().String()), nil
}

func (f *WsClientIDFactoryImpl) FromString(s string) (entity.WsClientPublicID, error) {
	if _, err := uuid.Parse(s); err != nil {
		return "", err
	}
	return entity.WsClientPublicID(s), nil
}
