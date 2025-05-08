package sqliteroomrepoimpl

import (
	"context"
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

func NewRoomRepositoryImpl(params *NewRoomRepositoryImplParams) repository.RoomRepository {
	if err := params.Validate(); err != nil {
		panic(err)
	}

	return &RoomRepositoryImpl{
		db: params.DB,
	}
}

func (r *RoomRepositoryImpl) SaveRoom(ctx context.Context, room *entity.Room) (entity.RoomID, error) {
	_, err := r.db.ExecContext(ctx, `INSERT INTO rooms (id, name) VALUES (?, ?)`, room.GetID(), room.GetName())
	if err != nil {
		return entity.RoomID(""), err
	}

	var roomID entity.RoomID
	err = r.db.GetContext(ctx, &roomID, `SELECT id FROM rooms WHERE id = ?`, room.GetID())
	if err != nil {
		return entity.RoomID(""), err
	}

	return roomID, nil
}

func (r *RoomRepositoryImpl) GetRoomByID(ctx context.Context, id entity.RoomID) (*entity.Room, error) {
	roomModel := model.RoomModel{}
	err := r.db.GetContext(ctx, &roomModel, `SELECT id, name FROM rooms WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}

	// ルームに所属しているユーザーを取得
	roomMembers := []model.RoomMemberModel{}
	err = r.db.SelectContext(ctx, &roomMembers, `SELECT user_id FROM room_members WHERE room_id = ?`, id)
	if err != nil {
		return nil, err
	}

	room := entity.NewRoom(entity.RoomParams{
		ID:       entity.RoomID(roomModel.ID.String()),
		Name:     roomModel.Name,
		Members:  make([]entity.UserID, len(roomMembers)),
	})
	return room, nil
}

func (r *RoomRepositoryImpl) GetAllRooms(ctx context.Context) ([]*entity.Room, error) {
	roomModels := []model.RoomModel{}
	err := r.db.SelectContext(ctx, &roomModels, `SELECT id, name FROM rooms`)
	if err != nil {
		return nil, err
	}

	rooms := make([]*entity.Room, len(roomModels))
	for i, roomModel := range roomModels {
		rooms[i] = entity.NewRoom(entity.RoomParams{
			ID:       entity.RoomID(roomModel.ID.String()),
			Name:     roomModel.Name,
			Members:  []entity.UserID{},
		})
	}

	return rooms, nil
}

func (r *RoomRepositoryImpl) GetUsersInRoom(ctx context.Context, roomID entity.RoomID) ([]*entity.User, error) {
	// まず中間テーブルから対象のuser_id一覧を取得
	var roomMembers []model.RoomMemberModel
	err := r.db.SelectContext(ctx, &roomMembers, `SELECT user_id FROM room_members WHERE room_id = ?`, roomID)
	if err != nil {
		return nil, err
	}

	// user_idだけ取り出す
	userIDs := make([]int, 0, len(roomMembers))
	for _, rm := range roomMembers {
		id, err := strconv.Atoi(rm.UserID.String())
		if err != nil {
			return nil, err
		}
		userIDs = append(userIDs, id)
	}

	if len(userIDs) == 0 {
		// メンバーがいなかった場合は空のスライスを返す
		return []*entity.User{}, nil
	}

	// IN句を使って一括でUser情報を取得（N+1回避！）
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
	_, err := r.db.ExecContext(ctx, `INSERT INTO room_members (room_id, user_id) VALUES (?, ?)`, roomID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *RoomRepositoryImpl) RemoveMemberFromRoom(ctx context.Context, roomID entity.RoomID, userID entity.UserID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM room_members WHERE room_id = ? AND user_id = ?`, roomID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *RoomRepositoryImpl) GetRoomByNameLike(ctx context.Context, name string) ([]*entity.Room, error) {
	roomModels := []model.RoomModel{}
	err := r.db.SelectContext(ctx, &roomModels, `SELECT id, name FROM rooms WHERE name LIKE ?`, "%"+name+"%")
	if err != nil {
		return nil, err
	}

	rooms := make([]*entity.Room, len(roomModels))
	for i, roomModel := range roomModels {
		rooms[i] = entity.NewRoom(entity.RoomParams{
			ID:       entity.RoomID(roomModel.ID.String()),
			Name:     roomModel.Name,
			Members:  []entity.UserID{},
		})
	}

	return rooms, nil
}

func (r *RoomRepositoryImpl) UpdateRoomName(ctx context.Context, roomID entity.RoomID, name string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE rooms SET name = ? WHERE id = ?`, name, roomID)
	if err != nil {
		return err
	}
	return nil
}

func (r *RoomRepositoryImpl) DeleteRoom(ctx context.Context, roomID entity.RoomID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM rooms WHERE id = ?`, roomID)
	if err != nil {
		return err
	}
	return nil
}

