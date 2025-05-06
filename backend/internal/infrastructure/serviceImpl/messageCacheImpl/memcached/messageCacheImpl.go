package memcachedmsgcacheimpl

import (
	"bytes"
	"encoding/gob"
	"errors"
	"time"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/domain/repository"
	"example.com/infrahandson/internal/domain/service"
	"github.com/bradfitz/gomemcache/memcache"
)

type messageCache struct {
	client  *memcache.Client
	msgRepo repository.MessageRepository
	limit   int
}

type NewMessageCacheServiceParams struct {
	MsgRepo repository.MessageRepository
	Client  *memcache.Client
}

func (p *NewMessageCacheServiceParams) Validate() error {
	if p.MsgRepo == nil {
		return errors.New("msgRepo is required")
	}
	if p.Client == nil {
		return errors.New("memcache Client is required")
	}
	return nil
}

// NewMessageCacheService は、MemCacheを使ったMessageCacheServiceの新しいインスタンスを返します。
func NewMessageCacheService(p *NewMessageCacheServiceParams) service.MessageCacheService {
	if err := p.Validate(); err != nil {
		panic(err)
	}
	return &messageCache{
		client:  p.Client,
		msgRepo: p.MsgRepo,
		limit:   service.DefaultRecentMessageLimit(),
	}
}

// Helper function to serialize messages to byte array for Memcached
func serializeMessages(messages []*entity.Message) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(messages); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Helper function to deserialize byte array back into messages
func deserializeMessages(data []byte) ([]*entity.Message, error) {
	var messages []*entity.Message
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&messages); err != nil {
		return nil, err
	}
	return messages, nil
}

func (m *messageCache) GetRecentMessages(roomID entity.RoomID) ([]*entity.Message, error) {
	// Memcachedでキャッシュを検索
	item, err := m.client.Get(string(roomID))
	if err == memcache.ErrCacheMiss {
		// キャッシュに存在しない場合は、リポジトリから取得
		messages, _, _, err := m.msgRepo.GetMessageHistoryInRoom(roomID, m.limit, time.Now())
		if err != nil {
			return nil, err
		}
		// 取得したメッセージをキャッシュに追加
		data, err := serializeMessages(messages)
		if err != nil {
			return nil, err
		}
		// Memcachedにセット（5分間キャッシュ）
		m.client.Set(&memcache.Item{
			Key:        string(roomID),
			Value:      data,
			Expiration: 5 * 60, // 5分
		})
		return messages, nil
	} else if err != nil {
		return nil, err
	}

	// キャッシュにあった場合
	messages, err := deserializeMessages(item.Value)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (m *messageCache) AddMessage(roomID entity.RoomID, message *entity.Message) error {
	// 既存のキャッシュを取得
	messages, err := m.GetRecentMessages(roomID)
	if err != nil {
		return err
	}

	// メッセージを追加
	messages = append([]*entity.Message{message}, messages...)

	// キャッシュがlimitを超えた場合、古いものを削除
	if len(messages) > m.limit {
		messages = messages[:m.limit]
	}

	// 更新したメッセージを再度キャッシュに保存
	data, err := serializeMessages(messages)
	if err != nil {
		return err
	}

	// Memcachedにセット（5分間キャッシュ）
	return m.client.Set(&memcache.Item{
		Key:        string(roomID),
		Value:      data,
		Expiration: 5 * 60, // 5分
	})
}
