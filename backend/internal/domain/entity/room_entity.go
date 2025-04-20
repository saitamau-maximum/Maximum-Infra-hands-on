package entity

type Room struct {
	id  RoomID
	public_id RoomPublicID
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
		public_id: params.PublicID,
		name: params.Name,
		members: params.Members,
	}
}

func (r *Room) GetPubID() RoomPublicID {
	// 部屋の公開IDを取得
	return r.public_id
}

func (r *Room) GetName() string {
	// 部屋の名前を取得
	return r.name
}

func (r *Room) GetMembers() []UserID {
	// 部屋のメンバーを取得
	return r.members
}