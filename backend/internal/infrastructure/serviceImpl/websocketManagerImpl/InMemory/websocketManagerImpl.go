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
	idByConn          map[service.WebSocketConnection]IDs
}

type IDs struct {
	UserID entity.UserID
	RoomID entity.RoomID
}

func NewInMemoryWebSocketManager() service.WebsocketManager {
	return &InMemoryWebSocketManager{
		connectionsByUser: make(map[entity.UserID]service.WebSocketConnection),
		connectionsByRoom: make(map[entity.RoomID]map[entity.UserID]service.WebSocketConnection),
		idByConn:			 make(map[service.WebSocketConnection]IDs),
	}
}

func (m *InMemoryWebSocketManager) Register(conn service.WebSocketConnection, userID entity.UserID, roomID entity.RoomID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	fmt.Println("\nRegistering connection for userID:", userID, "roomID:", roomID)
	// ユーザーごとの接続を登録
	m.connectionsByUser[userID] = conn

	// 部屋ごとの接続を登録
	if _, exists := m.connectionsByRoom[roomID]; !exists {
		m.connectionsByRoom[roomID] = make(map[entity.UserID]service.WebSocketConnection)
	}
	m.connectionsByRoom[roomID][userID] = conn
	// 接続とユーザーID、部屋IDのマッピングを保存
	m.idByConn[conn] = IDs{
		UserID: userID,
		RoomID: roomID,
	}
	return nil
}

func (m *InMemoryWebSocketManager) Unregister(conn service.WebSocketConnection) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	IDs, exists := m.idByConn[conn]
	if !exists {
		return errors.New("connection not found")
	}
	fmt.Println("\nUnregistering connection for userID:", IDs.UserID, "roomID:", IDs.RoomID)

	delete(m.idByConn, conn)
	delete(m.connectionsByUser, IDs.UserID)
	delete(m.connectionsByRoom[IDs.RoomID], IDs.UserID)
	
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

	for id, conn := range users {
		fmt.Println("Broadcasting to userID:", id)
		if err := conn.WriteMessage(msg); err != nil {
			return err
		}
	}
	return nil
}
