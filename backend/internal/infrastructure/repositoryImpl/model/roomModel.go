package model

import "github.com/google/uuid"

type RoomModel struct {
	ID   uuid.UUID `db:"id"`
	Name string `db:"name"`
}

type RoomMemberModel struct {
	ID     uuid.UUID `db:"id"`
	RoomID uuid.UUID `db:"room_id"`
	UserID uuid.UUID `db:"user_id"`
}
