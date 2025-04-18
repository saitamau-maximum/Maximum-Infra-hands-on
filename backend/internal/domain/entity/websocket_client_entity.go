package entity

type WebsocketClient struct {
	id       WebsocketClientID
	userID   UserID
	roomID   RoomID
}

type WebsocketClientParams struct {
	ID      WebsocketClientID
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

func (w *WebsocketClient) GetID() WebsocketClientID {
	return w.id
}