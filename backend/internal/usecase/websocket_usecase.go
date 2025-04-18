package usecase

import (
	"fmt"
	"time"

	"example.com/webrtc-practice/internal/domain/entity"
	"example.com/webrtc-practice/internal/domain/repository"
	"example.com/webrtc-practice/internal/domain/service"
	"example.com/webrtc-practice/internal/interface/factory"
)

type WebsocketUseCaseInterface interface {
	// 接続・参加処理
	ConnectUserToRoom(req ConnectUserToRoomRequest) error

	// メッセージ送信
	SendMessage(req SendMessageRequest) error

	// 切断処理
	DisconnectUser(req DisconnectUserRequest) error

	// メッセージ履歴取得
	GetMessageHistory(req GetMessageHistoryRequest) (GetMessageHistoryResponse, error)
}

type WebsocketUseCase struct {
	userRepo         repository.UserRepository
	roomRepo         repository.RoomRepository
	msgRepo          repository.MessageRepository
	wsClientRepo     repository.WebsocketClientRepository
	websocketManager service.WebsocketManager
	msgIDFactory     factory.MessageIDFactory
	clientIDFactory  factory.WebsocketClientIDFactory
}

func NewWebsocketUseCase(params NewWebsocketUseCaseParams) WebsocketUseCaseInterface {
	// Paramsのバリデーションを行う
	required := map[string]any{
		"UserRepo":         params.UserRepo,
		"RoomRepo":         params.RoomRepo,
		"MsgRepo":          params.MsgRepo,
		"WsClientRepo":     params.WsClientRepo,
		"WebsocketManager": params.WebsocketManager,
		"MsgIDFactory":     params.MsgIDFactory,
		"ClientIDFactory":  params.ClientIDFactory,
	}

	for name, f := range required {
		if f == nil {
			panic(fmt.Sprintf("%s is required", name))
		}
	}

	return &WebsocketUseCase{
		userRepo:         params.UserRepo,
		roomRepo:         params.RoomRepo,
		msgRepo:          params.MsgRepo,
		wsClientRepo:     params.WsClientRepo,
		websocketManager: params.WebsocketManager,
		msgIDFactory:     params.MsgIDFactory,
		clientIDFactory:  params.ClientIDFactory,
	}
}

// ConnectUserToRoom 接続・参加処理
func (w *WebsocketUseCase) ConnectUserToRoom(req ConnectUserToRoomRequest) error {
	user, err := w.userRepo.GetUserByID(req.UserID)
	if err != nil {
		return err
	}

	roomID, err := w.roomRepo.GetRoomIDByPublicID(req.PublicRoomID)
	if err != nil {
		return err
	}

	id, err := w.clientIDFactory.NewWebsocketClientID()
	if err != nil {
		return err
	}

	client := entity.NewWebsocketClient(entity.WebsocketClientParams{
		ID:     id,
		UserID: user.GetID(),
		RoomID: roomID,
	})

	err = w.wsClientRepo.CreateClient(client)
	if err != nil {
		return err
	}

	err = w.websocketManager.Register(req.Conn, req.UserID, roomID)
	if err != nil {
		return err
	}

	return nil
}

func (w *WebsocketUseCase) SendMessage(req SendMessageRequest) error {
	roomID, err := w.roomRepo.GetRoomIDByPublicID(req.RoomPublicID)
	if err != nil {
		return err
	}

	msgID, err := w.msgIDFactory.NewMessageID()
	if err != nil {
		return err
	}
	msg := entity.NewMessage(entity.MessageParams{
		ID:      msgID,
		RoomID:  roomID,
		UserID:  req.Sender,
		Content: req.Content,
		SentAt:  time.Now(),
	})

	err = w.websocketManager.BroadcastToRoom(roomID, msg)
	if err != nil {
		return err
	}

	return nil
}

func (w *WebsocketUseCase) DisconnectUser(req DisconnectUserRequest) error {
	conn, err := w.websocketManager.GetConnctionByUserID(req.UserID)
	if err != nil {
		return err
	}

	user, err := w.wsClientRepo.GetClientsByUserID(req.UserID)
	if err != nil {
		return err
	}

	err = w.websocketManager.Unregister(conn)
	if err != nil {
		return err
	}

	err = w.wsClientRepo.DeleteClient(user.GetID())
	if err != nil {
		return err
	}

	return nil
}

func (w *WebsocketUseCase) GetMessageHistory(req GetMessageHistoryRequest) (GetMessageHistoryResponse, error) {
	roomID, err := w.roomRepo.GetRoomIDByPublicID(req.PublicRoomID)
	if err != nil {
		return GetMessageHistoryResponse{nil}, err
	}

	messages, err := w.msgRepo.GetMessagesByRoomID(roomID)
	if err != nil {
		return GetMessageHistoryResponse{nil}, err
	}

	return GetMessageHistoryResponse{
		Messages: messages,
	}, nil
}

type NewWebsocketUseCaseParams struct {
	UserRepo         repository.UserRepository
	RoomRepo         repository.RoomRepository
	MsgRepo          repository.MessageRepository
	WsClientRepo     repository.WebsocketClientRepository
	WebsocketManager service.WebsocketManager
	MsgIDFactory     factory.MessageIDFactory
	ClientIDFactory  factory.WebsocketClientIDFactory
}

type ConnectUserToRoomRequest struct {
	UserID       entity.UserID
	PublicRoomID entity.RoomPublicID
	Conn         service.WebSocketConnection
}

type SendMessageRequest struct {
	RoomPublicID entity.RoomPublicID
	Sender       entity.UserID
	Content      string
}

type DisconnectUserRequest struct {
	UserID entity.UserID
}

type GetMessageHistoryRequest struct {
	PublicRoomID entity.RoomPublicID
}

type GetMessageHistoryResponse struct {
	Messages []*entity.Message
}
