// UseCase層をまとめた構造体
package usecase

import (
	"example.com/infrahandson/internal/usecase/messagecase"
	"example.com/infrahandson/internal/usecase/roomcase"
	"example.com/infrahandson/internal/usecase/usercase"
	"example.com/infrahandson/internal/usecase/websocketcase"
)

// UseCase はアプリケーションのユースケースをまとめた構造体です。
// DI 層で扱う他の層を組み立てるために必要な構造体を定義します。
type UseCase struct {
	UserUseCase      usercase.UserUseCaseInterface
	RoomUseCase      roomcase.RoomUseCaseInterface
	MessageUseCase   messagecase.MessageUseCaseInterface
	WebsocketUseCase websocketcase.WebsocketUseCaseInterface
}
