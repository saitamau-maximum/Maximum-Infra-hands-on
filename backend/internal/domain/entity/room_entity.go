package entity

type Room struct {
	id  RoomID
	publicID RoomPublicID
	name string
	members []UserPublicID
}

type RoomParams struct {
	ID 	  RoomID
	PublicID RoomPublicID
	Name string
	Members []UserPublicID
}

func NewRoom(params RoomParams) *Room {
	return &Room{
		id: params.ID,
		publicID: params.PublicID,
		name: params.Name,
		members: params.Members,
	}
}

func (r *Room) GetID() RoomID {
	// 部屋のIDを取得
	return r.id
}

func (r *Room) GetPubID() RoomPublicID {
	// 部屋の公開IDを取得
	return r.publicID
}

func (r *Room) GetName() string {
	// 部屋の名前を取得
	return r.name
}

func (r *Room) GetMembers() []UserPublicID {
	// 部屋のメンバーを取得
	return r.members
}
