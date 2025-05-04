package mysqlroomrepoimpl

import (
	"errors"
	"strconv"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/domain/repository"
	"example.com/infrahandson/internal/infrastructure/repositoryImpl/model"
	"github.com/jmoiron/sqlx"
)

type RoomRepositoryImpl struct {
	db *sqlx.DB
}

type NewRoomRepositoryImplParams struct {
	DB *sqlx.DB
}

func (p *NewRoomRepositoryImplParams) Validate() error {
	if p.DB == nil {
		return errors.New("db cannot be nil")
	}
	return nil
}

func NewRoomRepositoryImpl(p *NewRoomRepositoryImplParams) repository.RoomRepository {
	if err := p.Validate(); err != nil {
		panic(err)
	}

	return &RoomRepositoryImpl{
		db: p.DB,
	}
}

func (r *RoomRepositoryImpl) SaveRoom(room *entity.Room) (entity.RoomID, error) {
	_, err := r.db.Exec(`INSERT INTO rooms (id, name) VALUES (?, ?)`, room.GetID(), room.GetName())
	if err != nil {
			return entity.RoomID(""), err
	}
	return room.GetID(), nil
}

func (r *RoomRepositoryImpl) GetRoomByID(id entity.RoomID) (*entity.Room, error) {
	roomModel := model.RoomModel{}
	err := r.db.Get(&roomModel, `SELECT id, name FROM rooms WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}

	// ルームに所属しているユーザーを取得
	roomMembers := []model.RoomMemberModel{}
	err = r.db.Select(&roomMembers, `SELECT user_id FROM room_members WHERE room_id = ?`, id)
	if err != nil {
		return nil, err
	}

	room := entity.NewRoom(entity.RoomParams{
		ID:       entity.RoomID(roomModel.ID),
		Name:     roomModel.Name,
		Members:  make([]entity.UserID, len(roomMembers)),
	})

	return room, nil
}

func (r *RoomRepositoryImpl) GetRoomByPubID(pubID entity.RoomID) (*entity.Room, error) {
	roomModel := model.RoomModel{}
	err := r.db.Get(&roomModel, `SELECT id, name FROM rooms WHERE id = ?`, pubID)
	if err != nil {
		return nil, err
	}

	// ルームに所属しているユーザーを取得
	roomMembers := []model.RoomMemberModel{}
	err = r.db.Select(&roomMembers, `SELECT user_id FROM room_members WHERE room_id = ?`, roomModel.ID)
	if err != nil {
		return nil, err
	}

	room := entity.NewRoom(entity.RoomParams{
		ID:       entity.RoomID(roomModel.ID),
		Name:     roomModel.Name,
		Members:  make([]entity.UserID, len(roomMembers)),
	})

	return room, nil
}

func (r *RoomRepositoryImpl) GetAllRooms() ([]*entity.Room, error) {
	roomModels := []model.RoomModel{}
	err := r.db.Select(&roomModels, `SELECT id, name FROM rooms`)
	if err != nil {
		return nil, err
	}

	rooms := make([]*entity.Room, len(roomModels))
	for i, roomModel := range roomModels {
		rooms[i] = entity.NewRoom(entity.RoomParams{
			ID:       entity.RoomID(roomModel.ID),
			Name:     roomModel.Name,
			Members:  []entity.UserID{},
		})
	}

	return rooms, nil
}

func (r *RoomRepositoryImpl) GetUsersInRoom(roomID entity.RoomID) ([]*entity.User, error) {
	// まず中間テーブルから対象のuser_id一覧を取得
	var roomMembers []model.RoomMemberModel
	err := r.db.Select(&roomMembers, `SELECT user_id FROM room_members WHERE room_id = ?`, roomID)
	if err != nil {
		return nil, err
	}

	// user_idだけ取り出す
	userIDs := make([]int, 0, len(roomMembers))
	for _, rm := range roomMembers {
		id, err := strconv.Atoi(rm.UserID)
		if err != nil {
			return nil, err
		}
		userIDs = append(userIDs, id)
	}

	if len(userIDs) == 0 {
		// メンバーがいなかった場合は空のスライスを返す
		return []*entity.User{}, nil
	}

	// IN句を使って一括でUser情報を取得（N+1回避）
	query, args, err := sqlx.In(`SELECT id, name FROM users WHERE id IN (?)`, userIDs)
	if err != nil {
		return nil, err
	}
	query = r.db.Rebind(query)

	var userModels []model.UserModel
	err = r.db.Select(&userModels, query, args...)
	if err != nil {
		return nil, err
	}

	users := make([]*entity.User, 0, len(userModels))
	for _, um := range userModels {
		users = append(users, entity.NewUser(entity.UserParams{
			ID:         entity.UserID(um.ID),
			Name:       um.Name,
			Email:      um.Email,
			PasswdHash: um.PasswordHash,
			CreatedAt:  um.CreatedAt,
			UpdatedAt:  um.UpdatedAt,
		}))
	}

	return users, nil
}

func (r *RoomRepositoryImpl) AddMemberToRoom(roomID entity.RoomID, userID entity.UserID) error {
	_, err := r.db.Exec(`INSERT INTO room_members (room_id, user_id) VALUES (?, ?)`, roomID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *RoomRepositoryImpl) RemoveMemberFromRoom(roomID entity.RoomID, userID entity.UserID) error {
	_, err := r.db.Exec(`DELETE FROM room_members WHERE room_id = ? AND user_id = ?`, roomID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *RoomRepositoryImpl) GetRoomByNameLike(name string) ([]*entity.Room, error) {
	roomModels := []model.RoomModel{}
	// NOTE: FULLTEXT INDEXが前提
	err := r.db.Select(&roomModels, `
	SELECT id, name
	FROM rooms
	WHERE MATCH(name) AGAINST(? IN BOOLEAN MODE)
`, name)

	if err != nil {
		return nil, err
	}

	rooms := make([]*entity.Room, len(roomModels))
	for i, roomModel := range roomModels {
		rooms[i] = entity.NewRoom(entity.RoomParams{
			ID:       entity.RoomID(roomModel.ID),
			Name:     roomModel.Name,
			Members:  []entity.UserID{},
		})
	}

	return rooms, nil
}

func (r *RoomRepositoryImpl) UpdateRoomName(roomID entity.RoomID, name string) error {
	_, err := r.db.Exec(`UPDATE rooms SET name = ? WHERE id = ?`, name, roomID)
	if err != nil {
		return err
	}
	return nil
}

func (r *RoomRepositoryImpl) DeleteRoom(roomID entity.RoomID) error {
	_, err := r.db.Exec(`DELETE FROM rooms WHERE id = ?`, roomID)
	if err != nil {
		return err
	}
	return nil
}

func (r *RoomRepositoryImpl) GetRoomIDByID(id entity.RoomID) (entity.RoomID, error) {
	var roomID entity.RoomID
	err := r.db.Get(&roomID, `SELECT id FROM rooms WHERE id = ?`, id)
	if err != nil {
		return entity.RoomID(""), err
	}
	return roomID, nil
}

func (r *RoomRepositoryImpl) GetIDByRoomID(id entity.RoomID) (entity.RoomID, error) {
	roomID := entity.RoomID("")
	err := r.db.Get(&roomID, `SELECT id FROM rooms WHERE id = ?`, id)
	if err != nil {
		return entity.RoomID(""), err
	}
	if roomID == "" {
		return entity.RoomID(""), errors.New("room not found")
	}
	return roomID, nil
}
