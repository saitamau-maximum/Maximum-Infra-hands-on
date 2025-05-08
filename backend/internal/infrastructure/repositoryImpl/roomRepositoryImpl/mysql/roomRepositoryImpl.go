package mysqlroomrepoimpl

import (
	"context"
	"errors"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/domain/repository"
	"example.com/infrahandson/internal/infrastructure/repositoryImpl/model"
	"github.com/google/uuid"
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

func (r *RoomRepositoryImpl) SaveRoom(ctx context.Context, room *entity.Room) (entity.RoomID, error) {
	// RoomID -> UUID
	id := room.GetID()
	idUUID, err := id.RoomID2UUID()
	if err != nil {
		return entity.RoomID(""), err
	}
	// UUID -> BIN
	_, err = r.db.ExecContext(ctx, `INSERT INTO rooms (id, name) VALUES (UUID_TO_BIN(?), ?)`, idUUID, room.GetName())
	if err != nil {
		return entity.RoomID(""), err
	}
	return room.GetID(), nil
}

func (r *RoomRepositoryImpl) GetRoomByID(ctx context.Context, id entity.RoomID) (*entity.Room, error) {
	roomModel := model.RoomModel{}
	// RoomID -> UUID
	idUUID, err := id.RoomID2UUID()
	if err != nil {
		return nil, err
	}
	err = r.db.Get(&roomModel, `SELECT BIN_TO_UUID(id) AS id, name FROM rooms WHERE id = UUID_TO_BIN(?)`, idUUID)
	if err != nil {
		return nil, err
	}

	// ルームに所属しているユーザーを取得
	roomMembers := []model.RoomMemberModel{}
	err = r.db.SelectContext(ctx, &roomMembers, `SELECT BIN_TO_UUID(user_id) AS user_id FROM room_members WHERE room_id = UUID_TO_BIN(?)`, roomModel.ID)
	if err != nil {
		return nil, err
	}

	room := entity.NewRoom(entity.RoomParams{
		ID:      entity.RoomID(roomModel.ID.String()),
		Name:    roomModel.Name,
		Members: make([]entity.UserID, len(roomMembers)),
	})

	return room, nil
}

func (r *RoomRepositoryImpl) GetRoomByPubID(ctx context.Context, id entity.RoomID) (*entity.Room, error) {
	roomModel := model.RoomModel{}
	// RoomID -> UUID
	idUUID, err := id.RoomID2UUID()
	if err != nil {
		return nil, err
	}

	err = r.db.GetContext(ctx, &roomModel, `SELECT BIN_TO_UUID(id) AS id, name FROM rooms WHERE id = UUID_TO_BIN(?)`, idUUID)
	if err != nil {
		return nil, err
	}

	// ルームに所属しているユーザーを取得
	roomMembers := []model.RoomMemberModel{}
	err = r.db.SelectContext(ctx, &roomMembers, `SELECT BIN_TO_UUID(user_id) AS user_id FROM room_members WHERE room_id = UUID_TO_BIN(?)`, roomModel.ID)
	if err != nil {
		return nil, err
	}

	room := entity.NewRoom(entity.RoomParams{
		ID:      entity.RoomID(roomModel.ID.String()),
		Name:    roomModel.Name,
		Members: make([]entity.UserID, len(roomMembers)),
	})

	return room, nil
}

func (r *RoomRepositoryImpl) GetAllRooms(ctx context.Context) ([]*entity.Room, error) {
	roomModels := []model.RoomModel{}
	err := r.db.SelectContext(ctx, &roomModels, `SELECT BIN_TO_UUID(id) AS id, name FROM rooms`)
	if err != nil {
		return nil, err
	}

	rooms := make([]*entity.Room, len(roomModels))
	for i, roomModel := range roomModels {
		rooms[i] = entity.NewRoom(entity.RoomParams{
			ID:      entity.RoomID(roomModel.ID.String()),
			Name:    roomModel.Name,
			Members: []entity.UserID{},
		})
	}

	return rooms, nil
}

func (r *RoomRepositoryImpl) GetUsersInRoom(ctx context.Context, roomID entity.RoomID) ([]*entity.User, error) {
	// RoomID -> UUID
	roomIDUUID, err := roomID.RoomID2UUID()
	if err != nil {
		return nil, err
	}
	// まず中間テーブルから対象のuser_id一覧を取得
	var roomMembers []model.RoomMemberModel
	err = r.db.SelectContext(ctx, &roomMembers, `SELECT BIN_TO_UUID(user_id) AS user_id FROM room_members WHERE room_id = UUID_TO_BIN(?)`, roomIDUUID)
	if err != nil {
		return nil, err
	}

	// user_idだけ取り出す
	userIDs := make([]uuid.UUID, 0, len(roomMembers))
	for _, rm := range roomMembers {
		userIDs = append(userIDs, rm.UserID)
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
	err = r.db.SelectContext(ctx, &userModels, query, args...)
	if err != nil {
		return nil, err
	}

	users := make([]*entity.User, 0, len(userModels))
	for _, um := range userModels {
		users = append(users, entity.NewUser(entity.UserParams{
			ID:         entity.UserID(um.ID.String()),
			Name:       um.Name,
			Email:      um.Email,
			PasswdHash: um.PasswordHash,
			CreatedAt:  um.CreatedAt,
			UpdatedAt:  um.UpdatedAt,
		}))
	}

	return users, nil
}

func (r *RoomRepositoryImpl) AddMemberToRoom(ctx context.Context, roomID entity.RoomID, userID entity.UserID) error {
	// RoomID -> UUID
	roomIDUUID, err := roomID.RoomID2UUID()
	if err != nil {
		return err
	}
	// UserID -> UUID
	userIDUUID, err := userID.UserID2UUID()
	if err != nil {
		return err
	}
	// UUID -> BIN
	_, err = r.db.ExecContext(ctx, `INSERT INTO room_members (room_id, user_id) VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?))`, roomIDUUID, userIDUUID)
	if err != nil {
		return err
	}
	return nil
}

func (r *RoomRepositoryImpl) RemoveMemberFromRoom(ctx context.Context, roomID entity.RoomID, userID entity.UserID) error {
	// RoomID -> UUID
	roomIDUUID, err := roomID.RoomID2UUID()
	if err != nil {
		return err
	}
	// UserID -> UUID
	userIDUUID, err := userID.UserID2UUID()
	if err != nil {
		return err
	}
	// UUID -> BIN
	_, err = r.db.ExecContext(ctx, `DELETE FROM room_members WHERE room_id = UUID_TO_BIN(?) AND user_id = UUID_TO_BIN(?)`, roomIDUUID, userIDUUID)
	if err != nil {
		return err
	}
	return nil
}

func (r *RoomRepositoryImpl) GetRoomByNameLike(ctx context.Context, name string) ([]*entity.Room, error) {
	roomModels := []model.RoomModel{}
	// NOTE: FULLTEXT INDEXが前提
	err := r.db.SelectContext(ctx, &roomModels, `
	SELECT BIN_TO_UUID(id) AS id, name
	FROM rooms
	WHERE MATCH(name) AGAINST(? IN BOOLEAN MODE)
`, name)

	if err != nil {
		return nil, err
	}

	rooms := make([]*entity.Room, len(roomModels))
	for i, roomModel := range roomModels {
		rooms[i] = entity.NewRoom(entity.RoomParams{
			ID:      entity.RoomID(roomModel.ID.String()),
			Name:    roomModel.Name,
			Members: []entity.UserID{},
		})
	}

	return rooms, nil
}

func (r *RoomRepositoryImpl) UpdateRoomName(ctx context.Context, roomID entity.RoomID, name string) error {
	// RoomID -> UUID
	roomIDUUID, err := roomID.RoomID2UUID()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, `UPDATE rooms SET name = ? WHERE id = UUID_TO_BIN(?)`, name, roomIDUUID)
	if err != nil {
		return err
	}
	return nil
}

func (r *RoomRepositoryImpl) DeleteRoom(ctx context.Context, roomID entity.RoomID) error {
	// RoomID -> UUID
	roomIDUUID, err := roomID.RoomID2UUID()
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, `DELETE FROM rooms WHERE id = UUID_TO_BIN(?)`, roomIDUUID)
	if err != nil {
		return err
	}
	return nil
}
