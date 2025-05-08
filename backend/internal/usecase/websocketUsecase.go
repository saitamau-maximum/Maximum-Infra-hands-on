package usecase

import (
	"context"
	"errors"
	"time"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/domain/repository"
	"example.com/infrahandson/internal/domain/service"
	"example.com/infrahandson/internal/interface/factory"
)

type WebsocketUseCaseInterface interface {
	// 接続・参加処理
	ConnectUserToRoom(ctx context.Context, req ConnectUserToRoomRequest) error

	// メッセージ送信
	SendMessage(ctx context.Context, req SendMessageRequest) error

	// 切断処理
	DisconnectUser(ctx context.Context, req DisconnectUserRequest) error
}

type WebsocketUseCase struct {
	userRepo         repository.UserRepository
	roomRepo         repository.RoomRepository
	msgRepo          repository.MessageRepository
	msgCache         service.MessageCacheService
	wsClientRepo     repository.WebsocketClientRepository
	websocketManager service.WebsocketManager
	msgIDFactory     factory.MessageIDFactory
	clientIDFactory  factory.WsClientIDFactory
}

type NewWebsocketUseCaseParams struct {
	UserRepo         repository.UserRepository
	RoomRepo         repository.RoomRepository
	MsgRepo          repository.MessageRepository
	MsgCache         service.MessageCacheService
	WsClientRepo     repository.WebsocketClientRepository
	WebsocketManager service.WebsocketManager
	MsgIDFactory     factory.MessageIDFactory
	ClientIDFactory  factory.WsClientIDFactory
}

func (p *NewWebsocketUseCaseParams) Validate() error {
	if p.UserRepo == nil {
		return errors.New("UserRepo is required")
	}
	if p.RoomRepo == nil {
		return errors.New("RoomRepo is required")
	}
	if p.MsgRepo == nil {
		return errors.New("MsgRepo is required")
	}
	if p.MsgCache == nil {
		return errors.New("MsgCache is required")
	}
	if p.WsClientRepo == nil {
		return errors.New("WsClientRepo is required")
	}
	if p.WebsocketManager == nil {
		return errors.New("WebsocketManager is required")
	}
	if p.MsgIDFactory == nil {
		return errors.New("MsgIDFactory is required")
	}
	if p.ClientIDFactory == nil {
		return errors.New("ClientIDFactory is required")
	}
	return nil
}

func NewWebsocketUseCase(params NewWebsocketUseCaseParams) WebsocketUseCaseInterface {
	// Paramsのバリデーションを行う
	if err := params.Validate(); err != nil {
		panic(err)
	}

	return &WebsocketUseCase{
		userRepo:         params.UserRepo,
		roomRepo:         params.RoomRepo,
		msgRepo:          params.MsgRepo,
		msgCache:         params.MsgCache,
		wsClientRepo:     params.WsClientRepo,
		websocketManager: params.WebsocketManager,
		msgIDFactory:     params.MsgIDFactory,
		clientIDFactory:  params.ClientIDFactory,
	}
}

// ConnectUserToRoomRequest構造体: 接続・参加処理のリクエスト
type ConnectUserToRoomRequest struct {
	UserID entity.UserID
	RoomID entity.RoomID
	Conn   service.WebSocketConnection
}

// ConnectUserToRoom 接続・参加処理
func (w *WebsocketUseCase) ConnectUserToRoom(ctx context.Context, req ConnectUserToRoomRequest) error {
	user, err := w.userRepo.GetUserByID(ctx, req.UserID)
	if err != nil {
		return err
	}

	id, err := w.clientIDFactory.NewWsClientID()
	if err != nil {
		return err
	}

	client := entity.NewWebsocketClient(entity.WebsocketClientParams{
		ID:     id,
		UserID: user.GetID(),
		RoomID: req.RoomID,
	})

	err = w.wsClientRepo.CreateClient(ctx, client)
	if err != nil {
		return err
	}

	err = w.websocketManager.Register(ctx, req.Conn, req.UserID, req.RoomID)
	if err != nil {
		return err
	}

	return nil
}

// SendMessageRequest構造体: メッセージ送信リクエスト
type SendMessageRequest struct {
	RoomID  entity.RoomID
	Sender  entity.UserID
	Content string
}

// SendMessage メッセージ送信
func (w *WebsocketUseCase) SendMessage(ctx context.Context, req SendMessageRequest) error {
	id, err := w.msgIDFactory.NewMessageID()
	if err != nil {
		return err
	}

	msg := entity.NewMessage(entity.MessageParams{
		ID:      id,
		RoomID:  req.RoomID,
		UserID:  req.Sender,
		Content: req.Content,
		SentAt:  time.Now(),
	})

	if err := w.msgRepo.CreateMessage(ctx, msg); err != nil {
		return err
	}

	if err := w.msgCache.AddMessage(ctx, req.RoomID, msg); err != nil {
		return err
	}

	err = w.websocketManager.BroadcastToRoom(ctx, req.RoomID, msg)
	if err != nil {
		return err
	}

	return nil
}

// DisconnectUserRequest構造体: 切断処理リクエスト
type DisconnectUserRequest struct {
	UserID entity.UserID
}

// DisconnectUser 切断処理
func (w *WebsocketUseCase) DisconnectUser(ctx context.Context, req DisconnectUserRequest) error {
	conn, err := w.websocketManager.GetConnectionByUserID(ctx, req.UserID)
	if err != nil {
		return err
	}

	user, err := w.wsClientRepo.GetClientsByUserID(ctx, req.UserID)
	if err != nil {
		return err
	}

	err = w.websocketManager.Unregister(ctx, conn)
	if err != nil {
		return err
	}

	err = w.wsClientRepo.DeleteClient(ctx, user.GetID())
	if err != nil {
		return err
	}

	return nil
}

// GetMessageHistoryRequest構造体: メッセージ履歴取得リクエスト
type GetMessageHistoryRequest struct {
	RoomID entity.RoomID
}

// GetMessageHistoryResponse構造体: メッセージ履歴取得レスポンス
type GetMessageHistoryResponse struct {
	Messages []*entity.Message
}
