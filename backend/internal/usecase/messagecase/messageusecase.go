package messagecase

import (
	"example.com/infrahandson/internal/domain/repository"
	"example.com/infrahandson/internal/domain/service"
)

type MessageUseCase struct {
	msgRepo  repository.MessageRepository
	msgCache service.MessageCacheService
	roomRepo repository.RoomRepository
	userRepo repository.UserRepository
}
