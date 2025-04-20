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
}

type MessageUseCase struct {
	msgRepo  repository.MessageRepository
	roomRepo repository.RoomRepository
}

type NewMessageUseCaseParams struct {
	MsgRepo  repository.MessageRepository
	RoomRepo repository.RoomRepository
}

func (p *NewMessageUseCaseParams) Validate() error {
	if p.MsgRepo == nil {
		return errors.New("MsgRepo is required")
	}
	if p.RoomRepo == nil {
		return errors.New("RoomRepo is required")
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
