package entity

type Room struct {
	id  RoomID
	publick_id RoomPublicID
	name string
	members []UserID
}

type RoomParams struct {
	ID 	  RoomID
	PublicID RoomPublicID
	Name string
	Members []UserID
}

func NewRoom(params RoomParams) *Room {
	return &Room{
		id: params.ID,
		publick_id: params.PublicID,
		name: params.Name,
		members: params.Members,
	}
}