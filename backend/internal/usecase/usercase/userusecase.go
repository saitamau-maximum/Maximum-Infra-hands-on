package usercase

import (
	"example.com/infrahandson/internal/domain/repository"
	"example.com/infrahandson/internal/domain/service"
	"example.com/infrahandson/internal/interface/adapter"
	"example.com/infrahandson/internal/interface/factory"
)

type UserUseCase struct {
	userRepo      repository.UserRepository
	hasher        adapter.HasherAdapter
	tokenSvc      adapter.TokenServiceAdapter
	iconSvc       service.IconStoreService
	userIDFactory factory.UserIDFactory
}
