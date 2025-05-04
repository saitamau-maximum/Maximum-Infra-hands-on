package model

type RoomModel struct {
	ID       string    `db:"id"`
	PublicID string `db:"public_id"`
	Name     string `db:"name"`
}

type RoomMemberModel struct {
	ID     string    `db:"id"`
	RoomID string    `db:"room_id"`
	UserID string `db:"user_id"`
}
