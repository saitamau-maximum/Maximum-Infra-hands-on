package inmemorywsclientrepoimpl_test

import (
	"context"
	"testing"

	"example.com/infrahandson/internal/domain/entity"
	inmemorywsclientrepoimpl "example.com/infrahandson/internal/infrastructure/repositoryImpl/websocketClientRepositoryImpl/InMemory"
	"github.com/stretchr/testify/require"
)

func newTestClient(roomID entity.RoomID, userID entity.UserID) *entity.WebsocketClient {
	return entity.NewWebsocketClient(entity.WebsocketClientParams{
		ID: entity.WsClientID("test"),
		RoomID:   roomID,
		UserID:   userID,
	})
}

func TestInMemoryWebsocketClientRepository(t *testing.T) {
	repo := inmemorywsclientrepoimpl.NewInMemoryWebsocketClientRepository(inmemorywsclientrepoimpl.NewInMemoryWebsocketClientRepositoryParams{})

	roomID := entity.RoomID("test-room")
	userID := entity.UserID("test-user")
	client := newTestClient(roomID, userID)

	t.Run("CreateClient and GetClientByID success", func(t *testing.T) {
		err := repo.CreateClient(context.Background(), client)
		require.NoError(t, err)

		got, err := repo.GetClientByID(context.Background(), client.GetID())
		require.NoError(t, err)
		require.Equal(t, client, got)
	})

	t.Run("CreateClient duplicate", func(t *testing.T) {
		err := repo.CreateClient(context.Background(), client)
		require.Error(t, err) // 同じIDの登録はエラーになる
	})

	t.Run("GetClientsByRoomID success", func(t *testing.T) {
		clients, err := repo.GetClientsByRoomID(context.Background(), roomID)
		require.NoError(t, err)
		require.Len(t, clients, 1)
		require.Equal(t, client, clients[0])
	})

	t.Run("GetClientsByRoomID not found", func(t *testing.T) {
		otherRoomID := entity.RoomID("other-room")
		clients, err := repo.GetClientsByRoomID(context.Background(), otherRoomID)
		require.NoError(t, err)
		require.Nil(t, clients)
	})

	t.Run("GetClientsByUserID success", func(t *testing.T) {
		got, err := repo.GetClientsByUserID(context.Background(), userID)
		require.NoError(t, err)
		require.Equal(t, client, got)
	})

	t.Run("GetClientsByUserID not found", func(t *testing.T) {
		otherUserID := entity.UserID("other-user")
		got, err := repo.GetClientsByUserID(context.Background(), otherUserID)
		require.Error(t, err)
		require.Nil(t, got)
	})

	t.Run("DeleteClient success", func(t *testing.T) {
		err := repo.DeleteClient(context.Background(), client.GetID())
		require.NoError(t, err)

		// ちゃんと消えているか確認
		got, err := repo.GetClientByID(context.Background(), client.GetID())
		require.Error(t, err)
		require.Nil(t, got)

		roomClients, err := repo.GetClientsByRoomID(context.Background(), roomID)
		require.NoError(t, err)
		require.Nil(t, roomClients)

		gotUser, err := repo.GetClientsByUserID(context.Background(), userID)
		require.Error(t, err)
		require.Nil(t, gotUser)
	})

	t.Run("DeleteClient not found", func(t *testing.T) {
		err := repo.DeleteClient(context.Background(), client.GetID()) // 2回目の削除はエラー
		require.Error(t, err)
	})
}
