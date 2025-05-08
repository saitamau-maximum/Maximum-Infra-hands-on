package usecase

import (
	"context"
	"errors"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/domain/repository"
	"example.com/infrahandson/internal/interface/factory"
)

type RoomUseCaseInterface interface {
	CreateRoom(ctx context.Context, req CreateRoomRequest) (CreateRoomResponse, error)
	GetRoomByID(ctx context.Context, params GetRoomByIDRequest) (GetRoomByIDResponse, error)
	GetAllRooms(ctx context.Context) ([]*entity.Room, error)
	GetUsersInRoom(ctx context.Context, req GetUsersInRoomRequest) (GetUsersInRoomResponse, error)
	JoinRoom(ctx context.Context, req JoinRoomRequest) error
	LeaveRoom(ctx context.Context, req LeaveRoomRequest) error
	SearchRoom(ctx context.Context, req SearchRoomRequest) (SearchRoomResponse, error)
	UpdateRoomName(ctx context.Context, req UpdateRoomNameRequest) error
	DeleteRoom(ctx context.Context, req DeleteRoomRequest) error
}

// RoomUseCase構造体: 部屋に関するユースケースを管理
type RoomUseCase struct {
	roomRepo      repository.RoomRepository
	userRepo      repository.UserRepository
	roomIDFactory factory.RoomIDFactory
}

// NewRoomUseCaseParams構造体: RoomUseCaseの初期化に必要なパラメータ
type NewRoomUseCaseParams struct {
	RoomRepo      repository.RoomRepository
	UserRepo      repository.UserRepository
	RoomIDFactory factory.RoomIDFactory
}

func (p NewRoomUseCaseParams) Validate() error {
	if p.RoomRepo == nil {
		return errors.New("RoomRepo is required")
	}
	if p.UserRepo == nil {
		return errors.New("UserRepo is required")
	}
	if p.RoomIDFactory == nil {
		return errors.New("RoomIDFactory is required")
	}
	return nil
}

// NewRoomUseCase: RoomUseCaseのインスタンスを生成
func NewRoomUseCase(p NewRoomUseCaseParams) *RoomUseCase {
	if err := p.Validate(); err != nil {
		panic(err)
	}

	return &RoomUseCase{
		roomRepo:      p.RoomRepo,
		userRepo:      p.UserRepo,
		roomIDFactory: p.RoomIDFactory,
	}
}

// CreateRoomRequest構造体: 部屋作成リクエストのデータ
type CreateRoomRequest struct {
	Name string `json:"name"` // 部屋名
}

// CreateRoomResponse構造体: 部屋作成レスポンスのデータ
type CreateRoomResponse struct {
	Room *entity.Room `json:"room"` // 作成した部屋
}

// CreateRoom: 新しい部屋を作成
func (r *RoomUseCase) CreateRoom(ctx context.Context, req CreateRoomRequest) (CreateRoomResponse, error) {
	id, err := r.roomIDFactory.NewRoomID()
	if err != nil {
		return CreateRoomResponse{nil}, err
	}

	room := entity.NewRoom(entity.RoomParams{
		ID:      id,
		Name:    req.Name,
		Members: []entity.UserID{},
	})
	savedRoomID, err := r.roomRepo.SaveRoom(ctx, room)
	if err != nil {
		return CreateRoomResponse{nil}, err
	}
	res, err := r.roomRepo.GetRoomByID(ctx, savedRoomID)
	if err != nil {
		return CreateRoomResponse{nil}, err
	}

	return CreateRoomResponse{Room: res}, nil
}

// GetRoomByIDParams構造体: 公開IDで部屋を取得するためのパラメータ
type GetRoomByIDRequest struct {
	ID entity.RoomID `json:"id"`
}

// GetRoomByIDResponse構造体: 公開IDで部屋を取得した結果
type GetRoomByIDResponse struct {
	Room *entity.Room `json:"room"` // 取得した部屋
}

// GetRoomByID: 公開IDを使用して部屋を取得
func (r *RoomUseCase) GetRoomByID(ctx context.Context, req GetRoomByIDRequest) (GetRoomByIDResponse, error) {
	room, err := r.roomRepo.GetRoomByID(ctx, req.ID)
	if err != nil {
		return GetRoomByIDResponse{}, err
	}
	if room == nil {
		return GetRoomByIDResponse{}, errors.New("room not found")
	}
	return GetRoomByIDResponse{Room: room}, nil
}

// GetAllRooms: 全ての部屋を取得
func (r *RoomUseCase) GetAllRooms(ctx context.Context) ([]*entity.Room, error) {
	rooms, err := r.roomRepo.GetAllRooms(ctx)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

// GetUsersInRoomRequest構造体: 部屋内のユーザーを取得するリクエスト
type GetUsersInRoomRequest struct {
	ID entity.RoomID `json:"id"` // 公開ID
}

// GetUsersInRoomResponse構造体: 部屋内のユーザー取得結果
type GetUsersInRoomResponse struct {
	Users []*entity.User `json:"users"` // 部屋内のユーザーリスト
}

// GetUsersInRoom: 部屋内のユーザーを取得
func (r *RoomUseCase) GetUsersInRoom(ctx context.Context, req GetUsersInRoomRequest) (GetUsersInRoomResponse, error) {
	users, err := r.roomRepo.GetUsersInRoom(ctx, req.ID)
	if err != nil {
		return GetUsersInRoomResponse{}, err
	}

	return GetUsersInRoomResponse{Users: users}, nil
}

// JoinRoomRequest構造体: 部屋に参加するリクエスト
type JoinRoomRequest struct {
	RoomID entity.RoomID `json:"room_id"` // 部屋の公開ID
	UserID entity.UserID `json:"user_id"` // 参加するユーザー
}

// JoinRoom: 部屋にユーザーを参加させる
func (r *RoomUseCase) JoinRoom(ctx context.Context, req JoinRoomRequest) error {
	err := r.roomRepo.AddMemberToRoom(ctx, req.RoomID, req.UserID)
	if err != nil {
		return err
	}

	return nil
}

// LeaveRoomRequest構造体: 部屋から退出するリクエスト
type LeaveRoomRequest struct {
	RoomID entity.RoomID `json:"room_id"` // 部屋の公開ID
	UserID entity.UserID `json:"user_id"` // 退出するユーザーID
}

// LeaveRoom: 部屋からユーザーを退出させる
func (r *RoomUseCase) LeaveRoom(ctx context.Context, req LeaveRoomRequest) error {
	err := r.roomRepo.RemoveMemberFromRoom(ctx, req.RoomID, req.UserID)
	if err != nil {
		return err
	}

	return nil
}

// SearchRoomRequest構造体: 部屋を検索するリクエスト
type SearchRoomRequest struct {
	Name string `json:"name"` // 検索する部屋名
}

// SearchRoomResponse構造体: 部屋検索結果
type SearchRoomResponse struct {
	Rooms []*entity.Room `json:"rooms"` // 検索結果の部屋リスト
}

// SearchRoom: 部屋を名前で検索
func (r *RoomUseCase) SearchRoom(ctx context.Context, req SearchRoomRequest) (SearchRoomResponse, error) {
	rooms, err := r.roomRepo.GetRoomByNameLike(ctx, req.Name)
	if err != nil {
		return SearchRoomResponse{}, err
	}

	return SearchRoomResponse{Rooms: rooms}, nil
}

// UpdateRoomNameRequest構造体: 部屋名を更新するリクエスト
type UpdateRoomNameRequest struct {
	RoomID  entity.RoomID `json:"room_id"`
	NewName string        `json:"new_name"` // 新しい部屋名
}

// UpdateRoomName: 部屋名を更新
func (r *RoomUseCase) UpdateRoomName(ctx context.Context, req UpdateRoomNameRequest) error {
	err := r.roomRepo.UpdateRoomName(ctx, req.RoomID, req.NewName)
	if err != nil {
		return err
	}
	return nil
}

// DeleteRoomRequest構造体: 部屋を削除するリクエスト
type DeleteRoomRequest struct {
	RoomID entity.RoomID `json:"room_id"` // 部屋の公開ID
}

// DeleteRoom: 部屋を削除
func (r *RoomUseCase) DeleteRoom(ctx context.Context, req DeleteRoomRequest) error {
	err := r.roomRepo.DeleteRoom(ctx, req.RoomID)
	if err != nil {
		return err
	}
	return nil
}
