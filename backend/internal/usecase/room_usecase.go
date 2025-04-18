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
	roomRepo            repository.RoomRepository
	roomIDFactory       factory.RoomIDFactory
	roomPublicIDFactory factory.RoomPublicIDFactory
}

// NewRoomUseCaseParams構造体: RoomUseCaseの初期化に必要なパラメータ
type NewRoomUseCaseParams struct {
	RoomRepo            repository.RoomRepository
	RoomIDFactory       factory.RoomIDFactory
	RoomPublicIDFactory factory.RoomPublicIDFactory
}

// NewRoomUseCase: RoomUseCaseのインスタンスを生成
func NewRoomUseCase(p NewRoomUseCaseParams) *RoomUseCase {
	return &RoomUseCase{
		roomRepo:            p.RoomRepo,
		roomIDFactory:       p.RoomIDFactory,
		roomPublicIDFactory: p.RoomPublicIDFactory,
	}
}

// CreateRoom: 新しい部屋を作成
func (r *RoomUseCase) CreateRoom(req CreateRoomRequest) (CreateRoomResponse, error) {
	roomID, err := r.roomIDFactory.NewRoomID()
	if err != nil {
		return CreateRoomResponse{nil}, err
	}

	roomPublicID, err := r.roomPublicIDFactory.NewRoomPublicID()
	if err != nil {
		return CreateRoomResponse{nil}, err
	}

	room := entity.NewRoom(entity.RoomParams{
		ID:       roomID,
		Name:     req.Name,
		PublicID: roomPublicID,
		Members:  []entity.UserID{req.FirstUserID},
	})

	savedRoomID, err := r.roomRepo.SaveRoom(room)
	if err != nil {
		return CreateRoomResponse{nil}, err
	}

	return CreateRoomResponse{RoomID: &savedRoomID}, nil
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

// JoinRoom: 部屋にユーザーを参加させる
func (r *RoomUseCase) JoinRoom(req JoinRoomRequest) error {
	id, err := r.roomRepo.GetRoomIDByPublicID(req.RoomPublicID)
	if err != nil {
		return err
	}

	if id <= 0 {
		return errors.New("room not found")
	}

	err = r.roomRepo.AddMemberToRoom(id, req.User.GetID())
	if err != nil {
		return err
	}

	return nil
}

// LeaveRoom: 部屋からユーザーを退出させる
func (r *RoomUseCase) LeaveRoom(req LeaveRoomRequest) error {
	id, err := r.roomRepo.GetRoomIDByPublicID(req.RoomPublicID)
	if err != nil {
		return err
	}

	if id <= 0 {
		return errors.New("room not found")
	}

	err = r.roomRepo.RemoveMemberFromRoom(id, req.UserID)
	if err != nil {
		return err
	}

	return nil
}

// SearchRoom: 部屋を名前で検索
func (r *RoomUseCase) SearchRoom(req SearchRoomRequest) (SearchRoomResponse, error) {
	rooms, err := r.roomRepo.GetRoomByNameLike(req.Name)
	if err != nil {
		return SearchRoomResponse{}, err
	}

	return SearchRoomResponse{Rooms: rooms}, nil
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

// --- DTO構造体をここにまとめる ---

// CreateRoomRequest構造体: 部屋作成リクエストのデータ
type CreateRoomRequest struct {
	Name        string        `json:"name"`          // 部屋名
	FirstUserID entity.UserID `json:"first_user_id"` // 最初のユーザーID
}

// CreateRoomResponse構造体: 部屋作成レスポンスのデータ
type CreateRoomResponse struct {
	RoomID *entity.RoomID `json:"room_id"` // 作成された部屋のID
}

// IsRoomIDNil: RoomIDがnilかどうかを判定
func (CreateRoomRes *CreateRoomResponse) IsRoomIDNil() bool {
	return CreateRoomRes.RoomID == nil
}

// GetRoomID: RoomIDを取得（nilの場合はゼロ値を返す）
func (CreateRoomRes *CreateRoomResponse) GetRoomID() entity.RoomID {
	if CreateRoomRes.IsRoomIDNil() {
		return entity.RoomID(0) // 適切なゼロ値またはコンストラクタに置き換える
	}
	return *CreateRoomRes.RoomID
}

// GetRoomByPublicIDParams構造体: 公開IDで部屋を取得するためのパラメータ
type GetRoomByPublicIDParams struct {
	PublicID entity.RoomPublicID `json:"public_id"` // 公開ID
}

// GetRoomByPublicIDResponse構造体: 公開IDで部屋を取得した結果
type GetRoomByPublicIDResponse struct {
	Room *entity.Room `json:"room"` // 取得した部屋
}

// GetUsersInRoomRequest構造体: 部屋内のユーザーを取得するリクエスト
type GetUsersInRoomRequest struct {
	PublicID entity.RoomPublicID `json:"public_id"` // 公開ID
}

// GetUsersInRoomResponse構造体: 部屋内のユーザー取得結果
type GetUsersInRoomResponse struct {
	Users []*entity.User `json:"users"` // 部屋内のユーザーリスト
}

// JoinRoomRequest構造体: 部屋に参加するリクエスト
type JoinRoomRequest struct {
	RoomPublicID entity.RoomPublicID `json:"room_id"` // 部屋の公開ID
	User         *entity.User        `json:"user"`    // 参加するユーザー
}

// LeaveRoomRequest構造体: 部屋から退出するリクエスト
type LeaveRoomRequest struct {
	RoomPublicID entity.RoomPublicID `json:"room_id"` // 部屋の公開ID
	UserID       entity.UserID       `json:"user_id"` // 退出するユーザーID
}

// SearchRoomRequest構造体: 部屋を検索するリクエスト
type SearchRoomRequest struct {
	Name string `json:"name"` // 検索する部屋名
}

// SearchRoomResponse構造体: 部屋検索結果
type SearchRoomResponse struct {
	Rooms []*entity.Room `json:"rooms"` // 検索結果の部屋リスト
}

// UpdateRoomNameRequest構造体: 部屋名を更新するリクエスト
type UpdateRoomNameRequest struct {
	RoomPublicID entity.RoomPublicID `json:"room_id"`  // 部屋の公開ID
	NewName      string              `json:"new_name"` // 新しい部屋名
}

// DeleteRoomRequest構造体: 部屋を削除するリクエスト
type DeleteRoomRequest struct {
	RoomPublicID entity.RoomPublicID `json:"room_id"` // 部屋の公開ID
}
