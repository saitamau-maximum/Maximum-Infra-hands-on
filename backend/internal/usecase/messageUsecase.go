package usecase

import (
	"context"
	"errors"
	"time"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/domain/repository"
	"example.com/infrahandson/internal/domain/service"
)

type MessageUseCaseInterface interface {
	// メッセージ取得
	GetMessageHistoryInRoom(ctx context.Context, req GetMessageHistoryInRoomRequest) (GetMessageHistoryInRoomResponse, error)
}

type MessageUseCase struct {
	msgRepo  repository.MessageRepository
	msgCache service.MessageCacheService
	roomRepo repository.RoomRepository
	userRepo repository.UserRepository
}

type NewMessageUseCaseParams struct {
	MsgRepo  repository.MessageRepository
	MsgCache service.MessageCacheService
	RoomRepo repository.RoomRepository
	UserRepo repository.UserRepository
}

func (p *NewMessageUseCaseParams) Validate() error {
	if p.MsgRepo == nil {
		return errors.New("MsgRepo is required")
	}
	if p.MsgCache == nil {
		return errors.New("MsgCache is required")
	}
	if p.RoomRepo == nil {
		return errors.New("RoomRepo is required")
	}
	if p.UserRepo == nil {
		return errors.New("UserRepo is required")
	}
	return nil
}

func NewMessageUseCase(params NewMessageUseCaseParams) *MessageUseCase {
	if err := params.Validate(); err != nil {
		panic(err)
	}
	return &MessageUseCase{
		msgRepo:  params.MsgRepo,
		msgCache: params.MsgCache,
		roomRepo: params.RoomRepo,
		userRepo: params.UserRepo,
	}
}

type GetMessageHistoryInRoomRequest struct {
	RoomID       entity.RoomID
	Limit        int
	BeforeSentAt time.Time
}

type GetMessageHistoryInRoomResponse struct {
	Messages         []*entity.Message
	NextBeforeSentAt time.Time
	HasNext          bool
}

func (uc *MessageUseCase) GetMessageHistoryInRoom(ctx context.Context, req GetMessageHistoryInRoomRequest) (GetMessageHistoryInRoomResponse, error) {
	// まずはキャッシュからの取得を試みる
	messages, err := uc.msgCache.GetRecentMessages(ctx, req.RoomID)
	if err != nil {
		return GetMessageHistoryInRoomResponse{}, err
	}

	// キャッシュが空の場合、キャッシュ利用をスキップしてDBへ
	if len(messages) > 0 {
		// キャッシュの中で最も新しい・古いメッセージの時刻を調べる
		earliest := messages[0].GetSentAt()
		latest := messages[0].GetSentAt()
		for _, t := range messages[1:] {
			temp := t.GetSentAt()
			if temp.Before(earliest) {
				earliest = temp
			}
			if temp.After(latest) {
				latest = temp
			}
		}

		// キャッシュが使えるか判定（キャッシュの最新が BeforeSentAt よりも古い）
		if latest.Before(req.BeforeSentAt) {
			return GetMessageHistoryInRoomResponse{
				Messages:         messages,
				NextBeforeSentAt: earliest,
				HasNext:          len(messages) >= req.Limit,
			}, nil
		}
	}

	// メッセージ履歴を取得（キャッシュが使えない場合）
	messages, nextBeforeSentAt, hasNext, err := uc.msgRepo.GetMessageHistoryInRoom(
		ctx,
		req.RoomID,
		req.Limit,
		req.BeforeSentAt,
	)
	if err != nil {
		return GetMessageHistoryInRoomResponse{}, err
	}

	return GetMessageHistoryInRoomResponse{
		Messages:         messages,
		NextBeforeSentAt: nextBeforeSentAt,
		HasNext:          hasNext,
	}, nil
}

