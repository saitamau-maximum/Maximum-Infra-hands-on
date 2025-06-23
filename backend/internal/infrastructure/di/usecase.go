package di

import (
	"example.com/infrahandson/internal/domain/repository"
	"example.com/infrahandson/internal/domain/service"
	"example.com/infrahandson/internal/interface/adapter"
	"example.com/infrahandson/internal/interface/factory"
	"example.com/infrahandson/internal/usecase"
	"example.com/infrahandson/internal/usecase/roomcase"
	"example.com/infrahandson/internal/usecase/usercase"
	"example.com/infrahandson/internal/usecase/websocketcase"
)

type UseCaseDependency struct {
	Adapter adapter.Adapter
	Factory factory.Factory
	Repo    repository.Repository
	Svc     service.Service
}

// UseCaseInitialize はUseCase層の初期化を行います。
// 返り値 usecase.UseCase はUseCase層をまとめた構造体(詳細：internal/usecase/usecase.go)です。
func UseCaseInitialize(
	dep *UseCaseDependency,
) *usecase.UseCase {
	return &usecase.UseCase{
		UserUseCase: usercase.NewUserUseCase(usercase.NewUserUseCaseParams{
			UserRepo:      dep.Repo.UserRepository,
			Hasher:        dep.Adapter.HasherAdapter,
			TokenSvc:      dep.Adapter.TokenServiceAdapter,
			IconSvc:       dep.Svc.IconStoreService,
			UserIDFactory: dep.Factory.UserIDFactory,
		}),
		RoomUseCase: roomcase.NewRoomUseCase(roomcase.NewRoomUseCaseParams{
			RoomRepo:      dep.Repo.RoomRepository,
			RoomIDFactory: dep.Factory.RoomIDFactory,
			UserRepo:      dep.Repo.UserRepository,
		}),
		WebsocketUseCase: websocketcase.NewWebsocketUseCase(websocketcase.NewWebsocketUseCaseParams{
			UserRepo:         dep.Repo.UserRepository,
			RoomRepo:         dep.Repo.RoomRepository,
			MsgRepo:          dep.Repo.MessageRepository,
			MsgCache:         dep.Svc.MessageCacheService,
			WsClientRepo:     dep.Repo.WsClientRepository,
			WebsocketManager: dep.Svc.WebsocketManager,
			MsgIDFactory:     dep.Factory.MessageIDFactory,
			ClientIDFactory:  dep.Factory.WsClientIDFactory,
		}),
	}
}
