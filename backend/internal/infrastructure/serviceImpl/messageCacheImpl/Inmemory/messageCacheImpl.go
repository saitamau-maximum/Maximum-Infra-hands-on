package inmemorymsgcacheimpl

import (
	"container/list"
	"context"
	"errors"
	"sync"
	"time"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/domain/repository"
	"example.com/infrahandson/internal/domain/service"
)

type messageCache struct {
	mu      sync.RWMutex
	cache   map[entity.RoomID]*list.List // RoomID ごとに Message のリスト
	msgRepo repository.MessageRepository
	limit   int
}

type NewMessageCacheServiceParams struct {
	MsgRepo repository.MessageRepository
}

func (p *NewMessageCacheServiceParams) Validate() error {
	if p.MsgRepo == nil {
		return errors.New("MsgRepo is required")
	}
	return nil
}

// NewMessageCacheService は、MessageCacheService の新しいインスタンスを返します。
func NewMessageCacheService(p *NewMessageCacheServiceParams) service.MessageCacheService {
	if err := p.Validate(); err != nil {
		panic(err)
	}
	return &messageCache{
		cache: make(map[entity.RoomID]*list.List),
		msgRepo: p.MsgRepo,
		limit: service.DefaultRecentMessageLimit(),
	}
}

func (m *messageCache) GetRecentMessages(ctx context.Context, roomID entity.RoomID) ([]*entity.Message, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	lst, ok := m.cache[roomID]
	if !ok {
		// キャッシュに存在しない場合は、リポジトリから取得
		messages, _, _, err := m.msgRepo.GetMessageHistoryInRoom(ctx, roomID, m.limit, time.Now())
		if err != nil {
			return nil, err
		}
		// 取得したメッセージをキャッシュに追加
		lst = list.New()
		for _, msg := range messages {
			lst.PushBack(msg)
		}
		m.cache[roomID] = lst
	}

	messages := make([]*entity.Message, 0, lst.Len())
	for e := lst.Front(); e != nil; e = e.Next() {
		msg := e.Value.(*entity.Message)
		messages = append(messages, msg)
	}

	return messages, nil
}

func (m *messageCache) AddMessage(ctx context.Context, roomID entity.RoomID, message *entity.Message) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	lst, ok := m.cache[roomID]
	if !ok {
		lst = list.New()
		m.cache[roomID] = lst
	}

	// 末尾に追加（新しいメッセージ）
	lst.PushBack(message)

	// 古いものを削除（limit を超えた場合）
	for lst.Len() > m.limit {
		lst.Remove(lst.Front())
	}

	return nil
}
