package entity

type WebsocketClient struct {
	id       WsClientID
	userID   UserID
	roomID   RoomID
}

type WebsocketClientParams struct {
	ID      WsClientID
	UserID   UserID
	RoomID   RoomID
}

func NewWebsocketClient(params WebsocketClientParams) *WebsocketClient {
	return &WebsocketClient{
		id:     params.ID,
		userID: params.UserID,
		roomID: params.RoomID,
	}
}

func (w *WebsocketClient) GetID() WsClientID {
	return w.id
}

func (w *WebsocketClient) GetRoomID() RoomID {
	return w.roomID
}

func (w *WebsocketClient) GetUserID() UserID {
	return w.userID
}

func (w *WebsocketClient) SetID(id WsClientID) {
	w.id = id
}