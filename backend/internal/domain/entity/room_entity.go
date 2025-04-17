package entity

type Room struct {
	id  RoomID
	publick_id RoomPublicID
	name string
	members []UserID
}