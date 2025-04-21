package sqlite3

import (
	"errors"

	"example.com/webrtc-practice/internal/domain/entity"
	"example.com/webrtc-practice/internal/domain/repository"
	"github.com/jmoiron/sqlx"
)

type RoomRepositoryImpl struct {
	db *sqlx.DB
}

type NewRoomRepositoryParams struct {
	DB *sqlx.DB
}

func (p *NewRoomRepositoryParams) Validate() error {
	if p.DB == nil {
		return errors.New("DB is required")
	}
	return nil
}

func NewRoomRepository(params NewRoomRepositoryParams) repository.RoomRepository {
	if err := params.Validate(); err != nil {
		panic(err)
	}

	return &RoomRepositoryImpl{
		db: params.DB,
	}
}

type RoomModel struct {
	ID       int    `db:"id"`
	PublicID string `db:"public_id"`
	Name     string `db:"name"`
}

// SaveRoom inserts a new room into the database and returns its ID.
func (r *RoomRepositoryImpl) SaveRoom(room *entity.Room) (entity.RoomID, error) {
	query := `INSERT INTO rooms (public_id, name) VALUES (?, ?) RETURNING id`
	var roomID int
	err := r.db.QueryRow(query, string(room.GetPublicID()), room.GetName()).Scan(&roomID)
	if err != nil {
		return 0, err
	}
	return entity.RoomID(roomID), nil
}

// GetRoomByID retrieves a room by its ID.
func (r *RoomRepositoryImpl) GetRoomByID(id entity.RoomID) (*entity.Room, error) {
	query := `SELECT id, public_id, name FROM rooms WHERE id = ?`
	var model RoomModel
	err := r.db.Get(&model, query, int(id))
	if err != nil {
		return nil, err
	}
	return entity.NewRoom(entity.RoomParams{
		ID:       entity.RoomID(model.ID),
		PublicID: entity.RoomPublicID(model.PublicID),
		Name:     model.Name,
	}), nil
}

// GetAllRooms retrieves all rooms from the database.
func (r *RoomRepositoryImpl) GetAllRooms() ([]*entity.Room, error) {
	query := `SELECT id, public_id, name FROM rooms`
	var models []RoomModel
	err := r.db.Select(&models, query)
	if err != nil {
		return nil, err
	}

	var rooms []*entity.Room
	for _, model := range models {
		rooms = append(rooms, entity.NewRoom(entity.RoomParams{
			ID:       entity.RoomID(model.ID),
			PublicID: entity.RoomPublicID(model.PublicID),
			Name:     model.Name,
		}))
	}
	return rooms, nil
}

func (r *RoomRepositoryImpl) GetUsersInRoom(roomID entity.RoomID) ([]*entity.User, error) {
	query := `SELECT u.id, u.name FROM users u
			  INNER JOIN room_members rm ON u.id = rm.user_id
			  WHERE rm.room_id = ?`
	var userModels []struct {
		ID   string `db:"id"`
		Name string `db:"name"`
	}
	err := r.db.Select(&userModels, query, int(roomID))
	if err != nil {
		return nil, err
	}

	var users []*entity.User
	for _, model := range userModels {
		users = append(users, entity.NewUser(entity.UserParams{
			ID:   entity.UserID(model.ID),
			Name: model.Name,
		}))
	}
	return users, nil
}

func (r *RoomRepositoryImpl) AddMemberToRoom(roomID entity.RoomID, userID entity.UserID) error {
	query := `INSERT INTO room_members (room_id, user_id) VALUES (?, ?)`
	_, err := r.db.Exec(query, int(roomID), string(userID))
	return err
}

func (r *RoomRepositoryImpl) RemoveMemberFromRoom(roomID entity.RoomID, userID entity.UserID) error {
	query := `DELETE FROM room_members WHERE room_id = ? AND user_id = ?`
	_, err := r.db.Exec(query, int(roomID), string(userID))
	return err
}

// GetRoomByNameLike retrieves rooms whose names match the given pattern.
func (r *RoomRepositoryImpl) GetRoomByNameLike(name string) ([]*entity.Room, error) {
	query := `SELECT id, public_id, name FROM rooms WHERE name LIKE ?`
	var models []RoomModel
	err := r.db.Select(&models, query, "%"+name+"%")
	if err != nil {
		return nil, err
	}

	var rooms []*entity.Room
	for _, model := range models {
		rooms = append(rooms, entity.NewRoom(entity.RoomParams{
			ID:       entity.RoomID(model.ID),
			PublicID: entity.RoomPublicID(model.PublicID),
			Name:     model.Name,
		}))
	}
	return rooms, nil
}

// UpdateRoomName updates the name of a room.
func (r *RoomRepositoryImpl) UpdateRoomName(roomID entity.RoomID, newName string) error {
	query := `UPDATE rooms SET name = ? WHERE id = ?`
	_, err := r.db.Exec(query, newName, int(roomID))
	return err
}

// DeleteRoom deletes a room by its ID.
func (r *RoomRepositoryImpl) DeleteRoom(roomID entity.RoomID) error {
	query := `DELETE FROM rooms WHERE id = ?`
	_, err := r.db.Exec(query, int(roomID))
	return err
}

// GetRoomIDByPublicID retrieves a room's ID by its public ID.
func (r *RoomRepositoryImpl) GetRoomIDByPublicID(publicID entity.RoomPublicID) (entity.RoomID, error) {
	query := `SELECT id FROM rooms WHERE public_id = ?`
	var roomID int
	err := r.db.Get(&roomID, query, string(publicID))
	if err != nil {
		return 0, err
	}
	return entity.RoomID(roomID), nil
}

// GetPublicIDByRoomID retrieves a room's public ID by its ID.
func (r *RoomRepositoryImpl) GetPublicIDByRoomID(roomID entity.RoomID) (entity.RoomPublicID, error) {
	query := `SELECT public_id FROM rooms WHERE id = ?`
	var publicID string
	err := r.db.Get(&publicID, query, int(roomID))
	if err != nil {
		return "", err
	}
	return entity.RoomPublicID(publicID), nil
}
