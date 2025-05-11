package memcachedmsgcacheimpl

import (
	"bytes"
	"context"
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

// serializeMessages は、メッセージをDTO化してMemcached用にエンコードします。
func serializeMessages(_ context.Context, messages []*entity.Message) ([]byte, error) {
	dtos := make([]*MessageDTO, 0, len(messages))
	for _, m := range messages {
		dtos = append(dtos, fromEntityMessage(m))
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(dtos); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// deserializeMessages は、Memcachedのバイト列をDTOとしてデコードし、エンティティへ変換します。
func deserializeMessages(_ context.Context, data []byte) ([]*entity.Message, error) {
	var dtos []*MessageDTO
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&dtos); err != nil {
		return nil, err
	}

	messages := make([]*entity.Message, 0, len(dtos))
	for _, d := range dtos {
		msg, err := toEntityMessage(d)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

func (m *messageCache) GetRecentMessages(ctx context.Context, roomID entity.RoomID) ([]*entity.Message, error) {
	item, err := m.client.Get(string(roomID))
	if err == memcache.ErrCacheMiss {
		messages, _, _, err := m.msgRepo.GetMessageHistoryInRoom(ctx, roomID, m.limit, time.Now())
		if err != nil {
			return nil, err
		}
		data, err := serializeMessages(ctx, messages)

		if err != nil {
			return nil, err
		}
		m.client.Set(&memcache.Item{
			Key:        string(roomID),
			Value:      data,
			Expiration: 5 * 60,
		})
		return messages, nil
	} else if err != nil {
		return nil, err
	}

	messages, err := deserializeMessages(ctx, item.Value)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (m *messageCache) AddMessage(ctx context.Context, roomID entity.RoomID, message *entity.Message) error {
	messages, err := m.GetRecentMessages(ctx, roomID)
	if err != nil {
		return err
	}

	messages = append([]*entity.Message{message}, messages...)
	if len(messages) > m.limit {
		messages = messages[:m.limit]
	}

	data, err := serializeMessages(ctx, messages)
	if err != nil {
		return err
	}

	return m.client.Set(&memcache.Item{
		Key:        string(roomID),
		Value:      data,
		Expiration: 5 * 60,
	})
}
