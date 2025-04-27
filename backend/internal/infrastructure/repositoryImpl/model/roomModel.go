package model

type RoomModel struct {
	ID       int    `db:"id"`
	PublicID string `db:"public_id"`
	Name     string `db:"name"`
}

type RoomMemberModel struct {
	ID     int    `db:"id"`
	RoomID int    `db:"room_id"`
	UserID string `db:"user_id"`
}
