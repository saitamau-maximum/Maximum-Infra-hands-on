package usecase

import (
	"errors"
	"time"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/domain/repository"
	"example.com/infrahandson/internal/domain/service"
)

type MessageUseCaseInterface interface {
	// メッセージ取得
	GetMessageHistoryInRoom(req GetMessageHistoryInRoomRequest) (GetMessageHistoryInRoomResponse, error)
	// 内部のMessageを外部用に整形
	FormatMessage(msg *entity.Message) (FormatMessageResponse, error)
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

func (uc *MessageUseCase) GetMessageHistoryInRoom(req GetMessageHistoryInRoomRequest) (GetMessageHistoryInRoomResponse, error) {
	// メッセージ履歴を取得
	messages, nextBeforeSentAt, hasNext, err := uc.msgRepo.GetMessageHistoryInRoom(
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

type FormatMessageResponse struct {
	ID      string    `json:"id"`
	RoomID  string    `json:"room_id"`
	UserID  string    `json:"user_id"`
	Content string    `json:"content"`
	SentAt  time.Time `json:"sent_at"`
}

func (uc *MessageUseCase) FormatMessage(msg *entity.Message) (FormatMessageResponse, error) {
	if msg == nil {
		return FormatMessageResponse{}, errors.New("message is nil")
	}

	return FormatMessageResponse{
		ID:      string(msg.GetID()),
		RoomID:  string(msg.GetRoomID()),
		UserID:  string(msg.GetUserID()),
		Content: msg.GetContent(),
		SentAt:  msg.GetSentAt(),
	}, nil
}
