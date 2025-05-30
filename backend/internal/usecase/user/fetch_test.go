package user_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"example.com/infrahandson/internal/domain/entity"
	userUC "example.com/infrahandson/internal/usecase/user"
	mock_repository "example.com/infrahandson/test/mocks/domain/repository"
	mock_service "example.com/infrahandson/test/mocks/domain/service"
	mock_adapter "example.com/infrahandson/test/mocks/interface/adapter"
	mock_factory "example.com/infrahandson/test/mocks/interface/factory"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// 1.正常系
// 2.GetUserByID失敗

func TestGetUserByID(t *testing.T) {
	// Initialize the mock objects
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockHasher := mock_adapter.NewMockHasherAdapter(ctrl)
	mockTokenGenerator := mock_adapter.NewMockTokenServiceAdapter(ctrl)
	mockIconSvc := mock_service.NewMockIconStoreService(ctrl)
	mockUserIDFactory := mock_factory.NewMockUserIDFactory(ctrl)

	// Create the user use case with the mock dependencies
	userUseCase := userUC.NewUserUseCase(userUC.NewUserUseCaseParams{
		UserRepo:      mockUserRepo,
		Hasher:        mockHasher,
		TokenSvc:      mockTokenGenerator,
		IconSvc:       mockIconSvc,
		UserIDFactory: mockUserIDFactory,
	})

	t.Run("正常系", func(t *testing.T) {
		email := "test@mail.com"
		hashedPassword := "hashed_password"
		user := entity.NewUser(entity.UserParams{
			ID:         entity.UserID("user_id"),
			Name:       "John Doe",
			Email:      email,
			PasswdHash: hashedPassword,
			CreatedAt:  time.Now(),
			UpdatedAt:  nil,
		})
		mockUserRepo.EXPECT().GetUserByID(gomock.Any(), entity.UserID("user_id")).Return(user, nil)
		resp, err := userUseCase.GetUserByID(context.Background(), entity.UserID("user_id"))
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, user.GetID(), resp.GetID())
		assert.Equal(t, user.GetName(), resp.GetName())
		assert.Equal(t, user.GetEmail(), resp.GetEmail())
		assert.Equal(t, user.GetPasswdHash(), resp.GetPasswdHash())
		assert.NotNil(t, resp.GetCreatedAt())
	})

	t.Run("GetUserByID失敗", func(t *testing.T) {
		mockUserRepo.EXPECT().GetUserByID(gomock.Any(), entity.UserID("invalid_id")).Return(nil, errors.New("user not found"))
		resp, err := userUseCase.GetUserByID(context.Background(), entity.UserID("invalid_id"))
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}
