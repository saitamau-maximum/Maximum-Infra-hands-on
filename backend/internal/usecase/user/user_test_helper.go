package user

import (
	mock_repository "example.com/infrahandson/test/mocks/domain/repository"
	mock_service "example.com/infrahandson/test/mocks/domain/service"
	mock_adapter "example.com/infrahandson/test/mocks/interface/adapter"
	mock_factory "example.com/infrahandson/test/mocks/interface/factory"
	"go.uber.org/mock/gomock"
)

type mockDeps struct {
	UserRepo      *mock_repository.MockUserRepository
	Hasher        *mock_adapter.MockHasherAdapter
	TokenSvc      *mock_adapter.MockTokenServiceAdapter
	IconSvc       *mock_service.MockIconStoreService
	UserIDFactory *mock_factory.MockUserIDFactory
}

func NewTestUserUseCase(
	ctrl *gomock.Controller,
) (UserUseCaseInterface, mockDeps) {
	// モックの作成
	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockHasher := mock_adapter.NewMockHasherAdapter(ctrl)
	mockTokenSvc := mock_adapter.NewMockTokenServiceAdapter(ctrl)
	mockIconSvc := mock_service.NewMockIconStoreService(ctrl)
	mockUserIDFactory := mock_factory.NewMockUserIDFactory(ctrl)
	params := NewUserUseCaseParams{
		UserRepo:      mockUserRepo,
		Hasher:        mockHasher,
		TokenSvc:      mockTokenSvc,
		IconSvc:       mockIconSvc,
		UserIDFactory: mockUserIDFactory,
	}
	useCase := NewUserUseCase(params)

	return useCase, mockDeps{
		UserRepo:      mockUserRepo,
		Hasher:        mockHasher,
		TokenSvc:      mockTokenSvc,
		IconSvc:       mockIconSvc,
		UserIDFactory: mockUserIDFactory,
	}
}