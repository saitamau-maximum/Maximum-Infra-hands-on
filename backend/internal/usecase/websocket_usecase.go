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
	clientIDFactory  factory.WsClientIDFactory
}

type NewWebsocketUseCaseParams struct {
	UserRepo         repository.UserRepository
	RoomRepo         repository.RoomRepository
	MsgRepo          repository.MessageRepository
	WsClientRepo     repository.WebsocketClientRepository
	WebsocketManager service.WebsocketManager
	MsgIDFactory     factory.MessageIDFactory
	ClientIDFactory  factory.WsClientIDFactory
}
func (p *NewWebsocketUseCaseParams) Validate() error {
	if p.UserRepo == nil {
		return fmt.Errorf("UserRepo is required")
	}
	if p.RoomRepo == nil {
		return fmt.Errorf("RoomRepo is required")
	}
	if p.MsgRepo == nil {
		return fmt.Errorf("MsgRepo is required")
	}
	if p.WsClientRepo == nil {
		return fmt.Errorf("WsClientRepo is required")
	}
	if p.WebsocketManager == nil {
		return fmt.Errorf("WebsocketManager is required")
	}
	if p.MsgIDFactory == nil {
		return fmt.Errorf("MsgIDFactory is required")
	}
	if p.ClientIDFactory == nil {
		return fmt.Errorf("ClientIDFactory is required")
	}
	return nil
}

func NewWebsocketUseCase(params NewWebsocketUseCaseParams) WebsocketUseCaseInterface {
	// Paramsのバリデーションを行う
	if err := params.Validate(); err != nil {
		panic(fmt.Sprintf("Invalid parameters: %v", err))
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

// ConnectUserToRoomRequest構造体: 接続・参加処理のリクエスト
type ConnectUserToRoomRequest struct {
	UserID       entity.UserID
	PublicRoomID entity.RoomPublicID
	Conn         service.WebSocketConnection
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

	publicID, err := w.clientIDFactory.NewWsClientPublicID()
	if err != nil {
		return err
	}

	client := entity.NewWebsocketClient(entity.WebsocketClientParams{
		ID:     -1, // IDはDBに保存後に更新されるため、-1を指定
		PublicID: publicID,
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

// SendMessageRequest構造体: メッセージ送信リクエスト
type SendMessageRequest struct {
	RoomPublicID entity.RoomPublicID
	Sender       entity.UserID
	Content      string
}

// SendMessage メッセージ送信
func (w *WebsocketUseCase) SendMessage(req SendMessageRequest) error {
	roomID, err := w.roomRepo.GetRoomIDByPublicID(req.RoomPublicID)
	if err != nil {
		return err
	}

	msgPublicID, err := w.msgIDFactory.NewMessagePublicID()
	if err != nil {
		return err
	}
	msg := entity.NewMessage(entity.MessageParams{
		ID:      -1,
		PublicID: msgPublicID,
		RoomID:  roomID,
		UserID:  req.Sender,
		Content: req.Content,
		SentAt:  time.Now(),
	})

	if err := w.msgRepo.CreateMessage(msg); err != nil {
		return err
	}

	err = w.websocketManager.BroadcastToRoom(roomID, msg)
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
func (w *WebsocketUseCase) DisconnectUser(req DisconnectUserRequest) error {
	conn, err := w.websocketManager.GetConnectionByUserID(req.UserID)
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

// GetMessageHistoryRequest構造体: メッセージ履歴取得リクエスト
type GetMessageHistoryRequest struct {
	PublicRoomID entity.RoomPublicID
}

// GetMessageHistoryResponse構造体: メッセージ履歴取得レスポンス
type GetMessageHistoryResponse struct {
	Messages []*entity.Message
}

// GetMessageHistory メッセージ履歴取得
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


