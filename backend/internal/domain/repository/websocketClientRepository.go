package repository

import (
	"context"

	"example.com/infrahandson/internal/domain/entity"
)

type WebsocketClientRepository interface {
	CreateClient(ctx context.Context, client *entity.WebsocketClient) error
	DeleteClient(ctx context.Context, id entity.WsClientID) error
	GetClientByID(ctx context.Context, id entity.WsClientID) (*entity.WebsocketClient, error)
	GetClientsByRoomID(ctx context.Context, roomID entity.RoomID) ([]*entity.WebsocketClient, error)
	GetClientsByUserID(ctx context.Context, userID entity.UserID) (*entity.WebsocketClient, error)
}
