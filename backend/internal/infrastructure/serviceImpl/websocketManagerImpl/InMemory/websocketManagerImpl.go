package inmemorywsmanagerimpl

import (
	"errors"
	"fmt"
	"sync"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/domain/service"
)

type InMemoryWebSocketManager struct {
	mu                sync.RWMutex
	connectionsByUser map[entity.UserID]service.WebSocketConnection
	connectionsByRoom map[entity.RoomID]map[entity.UserID]service.WebSocketConnection
}

func NewInMemoryWebSocketManager() service.WebsocketManager {
	return &InMemoryWebSocketManager{
		connectionsByUser: make(map[entity.UserID]service.WebSocketConnection),
		connectionsByRoom: make(map[entity.RoomID]map[entity.UserID]service.WebSocketConnection),
	}
}

func (m *InMemoryWebSocketManager) Register(conn service.WebSocketConnection, userID entity.UserID, roomID entity.RoomID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// ユーザーごとの接続を登録
	m.connectionsByUser[userID] = conn

	// 部屋ごとの接続を登録
	if _, exists := m.connectionsByRoom[roomID]; !exists {
		m.connectionsByRoom[roomID] = make(map[entity.UserID]service.WebSocketConnection)
	}
	m.connectionsByRoom[roomID][userID] = conn
	return nil
}

func (m *InMemoryWebSocketManager) Unregister(conn service.WebSocketConnection) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// ユーザーIDを探して削除
	for userID, c := range m.connectionsByUser {
		if c == conn {
			delete(m.connectionsByUser, userID)
			break
		}
	}

	// 部屋IDを探して削除
	for roomID, users := range m.connectionsByRoom {
		for userID, c := range users {
			if c == conn {
				delete(users, userID)
				if len(users) == 0 {
					delete(m.connectionsByRoom, roomID)
				}
				break
			}
		}
	}
	return nil
}

func (m *InMemoryWebSocketManager) GetConnectionByUserID(userID entity.UserID) (service.WebSocketConnection, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	conn, exists := m.connectionsByUser[userID]
	if !exists {
		return nil, errors.New("connection not found for userID")
	}
	return conn, nil
}

func (m *InMemoryWebSocketManager) BroadcastToRoom(roomID entity.RoomID, msg *entity.Message) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	users, exists := m.connectionsByRoom[roomID]
	if !exists {
		return errors.New("room not found")
	}

	fmt.Println("Broadcast content:", msg.GetContent())

	for _, conn := range users {
		if err := conn.WriteMessage(msg); err != nil {
			return err
		}
	}
	return nil
}
