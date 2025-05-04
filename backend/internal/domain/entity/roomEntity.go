package entity

type Room struct {
	id      RoomID
	name    string
	members []UserID
}

type RoomParams struct {
	ID      RoomID
	Name    string
	Members []UserID
}

func NewRoom(params RoomParams) *Room {
	return &Room{
		id:      params.ID,
		name:    params.Name,
		members: params.Members,
	}
}

func (r *Room) GetID() RoomID {
	// 部屋のIDを取得
	return r.id
}

func (r *Room) GetName() string {
	// 部屋の名前を取得
	return r.name
}

func (r *Room) GetMembers() []UserID {
	// 部屋のメンバーを取得
	return r.members
}
