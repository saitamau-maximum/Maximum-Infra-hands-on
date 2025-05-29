// RoomのUseCaseの構造体
package room

import (
	"example.com/infrahandson/internal/domain/repository"
	"example.com/infrahandson/internal/interface/factory"
)

// RoomUseCase構造体: 部屋に関するユースケースを管理
type RoomUseCase struct {
	roomRepo      repository.RoomRepository
	userRepo      repository.UserRepository
	roomIDFactory factory.RoomIDFactory
}
