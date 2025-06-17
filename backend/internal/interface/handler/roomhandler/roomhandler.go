package roomhandler

import (
	"example.com/infrahandson/internal/interface/adapter"
	"example.com/infrahandson/internal/interface/factory"
	"example.com/infrahandson/internal/usecase/roomcase"
)

type RoomHandler struct {
	RoomUseCase   roomcase.RoomUseCaseInterface
	UserIDFactory factory.UserIDFactory
	RoomIDFactory factory.RoomIDFactory
	Logger        adapter.LoggerAdapter
}
