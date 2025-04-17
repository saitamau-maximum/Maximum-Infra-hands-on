package entity

type Room struct {
	id  RoomID
	name string
	members []UserID
}