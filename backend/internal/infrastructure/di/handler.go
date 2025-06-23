package di

import (
	"example.com/infrahandson/internal/interface/adapter"
	"example.com/infrahandson/internal/interface/factory"
	"example.com/infrahandson/internal/interface/handler"
	"example.com/infrahandson/internal/interface/handler/messagehandler"
	"example.com/infrahandson/internal/interface/handler/roomhandler"
	"example.com/infrahandson/internal/interface/handler/userhandler"
	"example.com/infrahandson/internal/interface/handler/websockethandler"
	"example.com/infrahandson/internal/usecase"
)

type HandlerInitializeParams struct {
	Adapter *adapter.Adapter
	Factory *factory.Factory
	UseCase *usecase.UseCase
}

func HandlerInitialize(
	params *HandlerInitializeParams,
) *handler.Handler {
	return &handler.Handler{
		UserHandler: userhandler.NewUserHandler(userhandler.NewUserHandlerParams{
			UserUseCase:   params.UseCase.UserUseCase,
			UserIDFactory: params.Factory.UserIDFactory,
			Logger:        params.Adapter.LoggerAdapter,
		}),
		RoomHandler: roomhandler.NewRoomHandler(roomhandler.NewRoomHandlerParams{
			RoomUseCase: params.UseCase.RoomUseCase,
			UserIDFactory: params.Factory.UserIDFactory,
			RoomIDFactory: params.Factory.RoomIDFactory,
			Logger:        params.Adapter.LoggerAdapter,
		}),
		WsHandler: websockethandler.NewWebSocketHandler(websockethandler.NewWebSocketHandlerParams{
			WsUseCase:   params.UseCase.WebsocketUseCase,
			WsUpgrader:  params.Adapter.Upgrader,
			WsConnFactory: params.Factory.WsConnFactory,
			UserIDFactory: params.Factory.UserIDFactory,
			RoomIDFactory: params.Factory.RoomIDFactory,
			Logger:        params.Adapter.LoggerAdapter,
		}),
		MsgHandler: messagehandler.NewMessageHandler(messagehandler.NewMessageHandlerParams{
			MsgUseCase: params.UseCase.MessageUseCase,
			Logger:     params.Adapter.LoggerAdapter,
		}),
	}
}
