package inmemorywsclientrepoimpl

import (
	"context"
	"errors"
	"sync"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/domain/repository"
)

type InMemoryWebsocketClientRepository struct {
	mu              sync.RWMutex
	clients         map[entity.WsClientID]*entity.WebsocketClient
	clientsByRoomID map[entity.RoomID]map[entity.WsClientID]*entity.WebsocketClient
	clientsByUserID map[entity.UserID]*entity.WebsocketClient
}

type NewInMemoryWebsocketClientRepositoryParams struct{}

func NewInMemoryWebsocketClientRepository(_ NewInMemoryWebsocketClientRepositoryParams) repository.WebsocketClientRepository {
	return &InMemoryWebsocketClientRepository{
		clients:         make(map[entity.WsClientID]*entity.WebsocketClient),
		clientsByRoomID: make(map[entity.RoomID]map[entity.WsClientID]*entity.WebsocketClient),
		clientsByUserID: make(map[entity.UserID]*entity.WebsocketClient),
	}
}

func (r *InMemoryWebsocketClientRepository) CreateClient(ctx context.Context, client *entity.WebsocketClient) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// IDを自動生成
	client.SetID(entity.WsClientID(client.GetID()))

	if _, exists := r.clients[client.GetID()]; exists {
		return errors.New("client already exists")
	}
	r.clients[client.GetID()] = client

	// RoomID側に登録
	if _, exists := r.clientsByRoomID[client.GetRoomID()]; !exists {
		r.clientsByRoomID[client.GetRoomID()] = make(map[entity.WsClientID]*entity.WebsocketClient)
	}
	r.clientsByRoomID[client.GetRoomID()][client.GetID()] = client

	// UserID側に登録
	r.clientsByUserID[client.GetUserID()] = client

	return nil
}

func (r *InMemoryWebsocketClientRepository) DeleteClient(ctx context.Context, id entity.WsClientID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	client, exists := r.clients[id]
	if !exists {
		return errors.New("client not found")
	}

	delete(r.clients, id)

	// RoomID側から削除
	if roomClients, exists := r.clientsByRoomID[client.GetRoomID()]; exists {
		delete(roomClients, id)
		if len(roomClients) == 0 {
			delete(r.clientsByRoomID, client.GetRoomID())
		}
	}

	// UserID側から削除
	delete(r.clientsByUserID, client.GetUserID())

	return nil
}

func (r *InMemoryWebsocketClientRepository) GetClientByID(ctx context.Context, id entity.WsClientID) (*entity.WebsocketClient, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	client, exists := r.clients[id]
	if !exists {
		return nil, errors.New("client not found")
	}
	return client, nil
}

func (r *InMemoryWebsocketClientRepository) GetClientsByRoomID(ctx context.Context, roomID entity.RoomID) ([]*entity.WebsocketClient, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	roomClients, exists := r.clientsByRoomID[roomID]
	if !exists {
		return nil, nil
	}

	var result []*entity.WebsocketClient
	for _, client := range roomClients {
		result = append(result, client)
	}
	return result, nil
}

func (r *InMemoryWebsocketClientRepository) GetClientsByUserID(ctx context.Context, userID entity.UserID) (*entity.WebsocketClient, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	client, exists := r.clientsByUserID[userID]
	if !exists {
		return nil, errors.New("client not found for user ID")
	}
	return client, nil
}
