package entity

type WebsocketClient struct {
	id       WsClientID
	publicID WsClientPublicID
	userID   UserID
	roomID   RoomID
}

type WebsocketClientParams struct {
	ID      WsClientID
	PublicID WsClientPublicID
	UserID   UserID
	RoomID   RoomID
}

func NewWebsocketClient(params WebsocketClientParams) *WebsocketClient {
	return &WebsocketClient{
		id:     params.ID,
		publicID: params.PublicID,
		userID: params.UserID,
		roomID: params.RoomID,
	}
}

func (w *WebsocketClient) GetID() WsClientID {
	return w.id
}