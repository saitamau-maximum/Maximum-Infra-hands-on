package usecase

import (
	"errors"
	"time"

	"example.com/webrtc-practice/internal/domain/entity"
	"example.com/webrtc-practice/internal/domain/repository"
)

type MessageUseCaseInterface interface {
	// メッセージ取得
	GetMessageHistoryInRoom(req GetMessageHistoryInRoomRequest) (GetMessageHistoryInRoomResponse, error)
	// 内部のMessageを外部用に整形
	FormatMessage(msg *entity.Message) (FormatMessageResponse, error)
}

type MessageUseCase struct {
	msgRepo  repository.MessageRepository
	roomRepo repository.RoomRepository
	userRepo repository.UserRepository
}

type NewMessageUseCaseParams struct {
	MsgRepo  repository.MessageRepository
	RoomRepo repository.RoomRepository
	UserRepo repository.UserRepository
}

func (p *NewMessageUseCaseParams) Validate() error {
	if p.MsgRepo == nil {
		return errors.New("MsgRepo is required")
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
		roomRepo: params.RoomRepo,
		userRepo: params.UserRepo,
	}
}

type GetMessageHistoryInRoomRequest struct {
	RoomPublicID entity.RoomPublicID
	Limit        int
	BeforeSentAt time.Time
}

type GetMessageHistoryInRoomResponse struct {
	Messages         []*entity.Message
	NextBeforeSentAt time.Time
	HasNext          bool
}

func (uc *MessageUseCase) GetMessageHistoryInRoom(req GetMessageHistoryInRoomRequest) (GetMessageHistoryInRoomResponse, error) {
	id, err := uc.roomRepo.GetRoomIDByPublicID(req.RoomPublicID)
	if err != nil {
		return GetMessageHistoryInRoomResponse{}, err
	}

	// メッセージ履歴を取得
	messages, nextBeforeSentAt, hasNext, err := uc.msgRepo.GetMessageHistoryInRoom(
		id,
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
	PublicID     string    `json:"id"`
	RoomPublicID string    `json:"room_id"`
	UserPublicID string    `json:"user_id"`
	Content      string    `json:"content"`
	SentAt       time.Time `json:"sent_at"`
}

func (uc *MessageUseCase) FormatMessage(msg *entity.Message) (FormatMessageResponse, error) {
	if msg == nil {
		return FormatMessageResponse{}, errors.New("message is nil")
	}

	roomPublicID, err := uc.roomRepo.GetPublicIDByRoomID(msg.GetRoomID())
	if err != nil {
		return FormatMessageResponse{}, err
	}

	userPublicID, err := uc.userRepo.GetPublicIDByID(msg.GetUserID())
	if err != nil {
		return FormatMessageResponse{}, err
	}

	return FormatMessageResponse{
		PublicID:     string(msg.GetPublicID()),
		RoomPublicID: string(roomPublicID),
		UserPublicID: string(userPublicID),
		Content:      msg.GetContent(),
		SentAt:       msg.GetSentAt(),
	}, nil
}
