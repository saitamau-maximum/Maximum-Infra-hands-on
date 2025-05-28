// WebsocketClientの永続化に関するインターフェース
// 具体実装は/infrastructure/repositoryImpl
package repository

import (
	"context"

	"example.com/infrahandson/internal/domain/entity"
)

// WebsocketClientRepositoryはWebSocketクライアントの永続化操作を定義するインターフェースです。
type WebsocketClientRepository interface {
	// CreateClientは新しいWebSocketクライアントを作成・保存します。
	CreateClient(ctx context.Context, client *entity.WebsocketClient) error

	// DeleteClientは指定されたIDのWebSocketクライアントを削除します。
	DeleteClient(ctx context.Context, id entity.WsClientID) error

	// GetClientByIDはクライアントIDに対応するWebSocketクライアントを取得します。
	GetClientByID(ctx context.Context, id entity.WsClientID) (*entity.WebsocketClient, error)

	// GetClientsByRoomIDは指定された部屋IDに接続しているすべてのWebSocketクライアントを取得します。
	GetClientsByRoomID(ctx context.Context, roomID entity.RoomID) ([]*entity.WebsocketClient, error)

	// GetClientsByUserIDは指定されたユーザーIDに紐づくWebSocketクライアントを取得します。
	GetClientsByUserID(ctx context.Context, userID entity.UserID) (*entity.WebsocketClient, error)
}

