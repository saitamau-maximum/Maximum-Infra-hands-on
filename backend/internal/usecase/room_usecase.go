package usecase

import (
	"errors"

	"example.com/webrtc-practice/internal/domain/entity"
	"example.com/webrtc-practice/internal/domain/repository"
	"example.com/webrtc-practice/internal/interface/factory"
)

type RoomUseCaseInterface interface {
	CreateRoom(req CreateRoomRequest) (CreateRoomResponse, error)
	GetRoomByPublicID(params GetRoomByPublicIDParams) (GetRoomByPublicIDResponse, error)
	GetAllRooms() ([]*entity.Room, error)
	GetUsersInRoom(req GetUsersInRoomRequest) (GetUsersInRoomResponse, error)
	JoinRoom(req JoinRoomRequest) error
	LeaveRoom(req LeaveRoomRequest) error
	SearchRoom(req SearchRoomRequest) (SearchRoomResponse, error)
	UpdateRoomName(req UpdateRoomNameRequest) error
	DeleteRoom(req DeleteRoomRequest) error
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
func (r *RoomUseCase) CreateRoom(req CreateRoomRequest) (CreateRoomResponse, error) {
	roomPublicID, err := r.roomIDFactory.NewRoomPublicID()
	if err != nil {
		return CreateRoomResponse{nil}, err
	}

	room := entity.NewRoom(entity.RoomParams{
		ID:       -1, // IDはDBに保存後に更新されるため、-1を指定
		PublicID: roomPublicID,
		Name:     req.Name,
		Members:  []entity.UserPublicID{},
	})

	savedRoomID, err := r.roomRepo.SaveRoom(room)
	if err != nil {
		return CreateRoomResponse{nil}, err
	}

	res, err := r.roomRepo.GetRoomByID(savedRoomID)
	if err != nil {
		return CreateRoomResponse{nil}, err
	}

	return CreateRoomResponse{Room: res}, nil
}

// GetRoomByPublicIDParams構造体: 公開IDで部屋を取得するためのパラメータ
type GetRoomByPublicIDParams struct {
	PublicID entity.RoomPublicID `json:"public_id"` // 公開ID
}

// GetRoomByPublicIDResponse構造体: 公開IDで部屋を取得した結果
type GetRoomByPublicIDResponse struct {
	Room *entity.Room `json:"room"` // 取得した部屋
}

// GetRoomByPublicID: 公開IDを使用して部屋を取得
func (r *RoomUseCase) GetRoomByPublicID(params GetRoomByPublicIDParams) (GetRoomByPublicIDResponse, error) {
	roomID, err := r.roomRepo.GetRoomIDByPublicID(params.PublicID)
	if err != nil {
		return GetRoomByPublicIDResponse{}, err
	}

	if roomID <= 0 {
		return GetRoomByPublicIDResponse{}, errors.New("room not found")
	}

	room, err := r.roomRepo.GetRoomByID(roomID)
	if err != nil {
		return GetRoomByPublicIDResponse{}, err
	}
	if room == nil {
		return GetRoomByPublicIDResponse{}, errors.New("room not found")
	}
	return GetRoomByPublicIDResponse{Room: room}, nil
}

// GetAllRooms: 全ての部屋を取得
func (r *RoomUseCase) GetAllRooms() ([]*entity.Room, error) {
	rooms, err := r.roomRepo.GetAllRooms()
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

// GetUsersInRoomRequest構造体: 部屋内のユーザーを取得するリクエスト
type GetUsersInRoomRequest struct {
	PublicID entity.RoomPublicID `json:"public_id"` // 公開ID
}

// GetUsersInRoomResponse構造体: 部屋内のユーザー取得結果
type GetUsersInRoomResponse struct {
	Users []*entity.User `json:"users"` // 部屋内のユーザーリスト
}

// GetUsersInRoom: 部屋内のユーザーを取得
func (r *RoomUseCase) GetUsersInRoom(req GetUsersInRoomRequest) (GetUsersInRoomResponse, error) {
	roomID, err := r.roomRepo.GetRoomIDByPublicID(req.PublicID)
	if err != nil {
		return GetUsersInRoomResponse{}, err
	}

	if roomID <= 0 {
		return GetUsersInRoomResponse{}, errors.New("room not found")
	}

	users, err := r.roomRepo.GetUsersInRoom(roomID)
	if err != nil {
		return GetUsersInRoomResponse{}, err
	}

	return GetUsersInRoomResponse{Users: users}, nil
}

// JoinRoomRequest構造体: 部屋に参加するリクエスト
type JoinRoomRequest struct {
	RoomPublicID entity.RoomPublicID `json:"room_id"` // 部屋の公開ID
	UserPublicID entity.UserPublicID `json:"user_id"` // 参加するユーザー
}

// JoinRoom: 部屋にユーザーを参加させる
func (r *RoomUseCase) JoinRoom(req JoinRoomRequest) error {
	roomPublicID, err := r.roomRepo.GetRoomIDByPublicID(req.RoomPublicID)
	if err != nil {
		return err
	}

	userPublicID, err := r.userRepo.GetIDByPublicID(req.UserPublicID)
	if err != nil {
		return err
	}

	if roomPublicID <= 0 {
		return errors.New("room not found")
	}
	if userPublicID <= 0 {
		return errors.New("user not found")
	}

	err = r.roomRepo.AddMemberToRoom(roomPublicID, userPublicID)
	if err != nil {
		return err
	}

	return nil
}

// LeaveRoomRequest構造体: 部屋から退出するリクエスト
type LeaveRoomRequest struct {
	RoomPublicID entity.RoomPublicID `json:"room_id"` // 部屋の公開ID
	UserPublicID entity.UserPublicID `json:"user_id"` // 退出するユーザーID
}

// LeaveRoom: 部屋からユーザーを退出させる
func (r *RoomUseCase) LeaveRoom(req LeaveRoomRequest) error {
	roomPublicID, err := r.roomRepo.GetRoomIDByPublicID(req.RoomPublicID)
	if err != nil {
		return err
	}

	userPublicID, err := r.userRepo.GetIDByPublicID(req.UserPublicID)
	if err != nil {
		return err
	}

	if userPublicID <= 0 {
		return errors.New("room not found")
	}

	if roomPublicID <= 0 {
		return errors.New("user not found")
	}

	err = r.roomRepo.RemoveMemberFromRoom(roomPublicID, userPublicID)
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
func (r *RoomUseCase) SearchRoom(req SearchRoomRequest) (SearchRoomResponse, error) {
	rooms, err := r.roomRepo.GetRoomByNameLike(req.Name)
	if err != nil {
		return SearchRoomResponse{}, err
	}

	return SearchRoomResponse{Rooms: rooms}, nil
}

// UpdateRoomNameRequest構造体: 部屋名を更新するリクエスト
type UpdateRoomNameRequest struct {
	RoomPublicID entity.RoomPublicID `json:"room_id"`  // 部屋の公開ID
	NewName      string              `json:"new_name"` // 新しい部屋名
}

// UpdateRoomName: 部屋名を更新
func (r *RoomUseCase) UpdateRoomName(req UpdateRoomNameRequest) error {
	roomID, err := r.roomRepo.GetRoomIDByPublicID(req.RoomPublicID)
	if err != nil {
		return err
	}

	if roomID <= 0 {
		return errors.New("room not found")
	}
	err = r.roomRepo.UpdateRoomName(roomID, req.NewName)
	if err != nil {
		return err
	}
	return nil
}

// DeleteRoomRequest構造体: 部屋を削除するリクエスト
type DeleteRoomRequest struct {
	RoomPublicID entity.RoomPublicID `json:"room_id"` // 部屋の公開ID
}

// DeleteRoom: 部屋を削除
func (r *RoomUseCase) DeleteRoom(req DeleteRoomRequest) error {
	roomID, err := r.roomRepo.GetRoomIDByPublicID(req.RoomPublicID)
	if err != nil {
		return err
	}

	if roomID <= 0 {
		return errors.New("room not found")
	}
	err = r.roomRepo.DeleteRoom(roomID)
	if err != nil {
		return err
	}
	return nil
}
