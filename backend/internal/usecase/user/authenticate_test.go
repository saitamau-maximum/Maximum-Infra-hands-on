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
// 2.GetUserByEmail失敗
// 3.ComparePassword失敗
// 4.パスワード不一致
// 5.トークン生成失敗
// 6.トークンの有効期限取得失敗

func TestAuthenticateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// == Mock dependencies ==
	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockHasher := mock_adapter.NewMockHasherAdapter(ctrl)
	mockTokenSvc := mock_adapter.NewMockTokenServiceAdapter(ctrl)
	mockIconSvc := mock_service.NewMockIconStoreService(ctrl)
	mockUserIDFactory := mock_factory.NewMockUserIDFactory(ctrl)

	params := userUC.NewUserUseCaseParams{
		UserRepo:      mockUserRepo,
		Hasher:        mockHasher,
		TokenSvc:      mockTokenSvc,
		IconSvc:       mockIconSvc,
		UserIDFactory: mockUserIDFactory,
	}

	userUseCase := userUC.NewUserUseCase(params)

	t.Run("正常系", func(t *testing.T) {
		email := "test@email.com"
		password := "password123"
		hashedPassword := "hashed_password"

		parms := entity.UserParams{
			ID:         entity.UserID("user_id"),
			Name:       "John Doe",
			Email:      email,
			PasswdHash: hashedPassword,
			CreatedAt:  time.Now(),
			UpdatedAt:  nil,
		}

		user := entity.NewUser(parms)
		token := "generated_token"

		req := userUC.AuthenticateUserRequest{
			Email:    email,
			Password: password,
		}

		mockUserRepo.EXPECT().GetUserByEmail(context.Background(), email).Return(user, nil)
		mockHasher.EXPECT().ComparePassword(hashedPassword, password).Return(true, nil)
		mockTokenSvc.EXPECT().GenerateToken(user.GetID()).Return(token, nil)
		mockTokenSvc.EXPECT().GetExpireAt(token).Return(1, nil)

		response, err := userUseCase.AuthenticateUser(context.Background(), req)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, response.GetToken(), token)
	})

	t.Run("GetUserByEmail失敗", func(t *testing.T) {
		email := "test@mail.com"
		password := "password123"
		req := userUC.AuthenticateUserRequest{
			Email:    email,
			Password: password,
		}
		mockUserRepo.EXPECT().GetUserByEmail(context.Background(), email).Return(nil, errors.New("user not found"))
		response, err := userUseCase.AuthenticateUser(context.Background(), req)
		assert.Error(t, err)
		assert.NotNil(t, response)
		assert.True(t, response.IsTokenNil())
	})

	t.Run("ComparePassword失敗", func(t *testing.T) {
		email := "test@mail.com"
		password := "password123"
		hashedPassword := "hashed_password"
		req := userUC.AuthenticateUserRequest{
			Email:    email,
			Password: password,
		}
		mockUserRepo.EXPECT().GetUserByEmail(context.Background(), email).Return(entity.NewUser(entity.UserParams{
			ID:         entity.UserID("user_id"),
			Name:       "John Doe",
			Email:      email,
			PasswdHash: hashedPassword,
			CreatedAt:  time.Now(),
			UpdatedAt:  nil,
		}), nil)

		mockHasher.EXPECT().ComparePassword(hashedPassword, password).Return(false, errors.New("hash compare failed"))
		response, err := userUseCase.AuthenticateUser(context.Background(), req)
		assert.Error(t, err)
		assert.NotNil(t, response)
		assert.True(t, response.IsTokenNil())
	})

	t.Run("パスワード不一致", func(t *testing.T) {
		email := "test@mail.om"
		password := "password123"
		hashedPassword := "hashed_password"
		req := userUC.AuthenticateUserRequest{
			Email:    email,
			Password: password,
		}

		mockUserRepo.EXPECT().GetUserByEmail(context.Background(), email).Return(entity.NewUser(entity.UserParams{
			ID:         entity.UserID("user_id"),
			Name:       "John Doe",
			Email:      email,
			PasswdHash: hashedPassword,
			CreatedAt:  time.Now(),
			UpdatedAt:  nil,
		}), nil)
		mockHasher.EXPECT().ComparePassword(hashedPassword, password).Return(false, nil)
		response, err := userUseCase.AuthenticateUser(context.Background(), req)
		assert.Error(t, err)
		assert.NotNil(t, response)
		assert.True(t, response.IsTokenNil())
	})

	t.Run("トークン生成失敗", func(t *testing.T) {
		email := "test@mail.com"
		password := "password123"
		hashedPassword := "hashed_password"
		req := userUC.AuthenticateUserRequest{
			Email:    email,
			Password: password,
		}

		mockUserRepo.EXPECT().GetUserByEmail(context.Background(), email).Return(entity.NewUser(entity.UserParams{
			ID:         entity.UserID("user_id"),
			Name:       "John Doe",
			Email:      email,
			PasswdHash: hashedPassword,
			CreatedAt:  time.Now(),
			UpdatedAt:  nil,
		}), nil)

		mockHasher.EXPECT().ComparePassword(hashedPassword, password).Return(true, nil)
		mockTokenSvc.EXPECT().GenerateToken(gomock.Any()).Return("", errors.New("token generation failed"))
		response, err := userUseCase.AuthenticateUser(context.Background(), req)
		assert.Error(t, err)
		assert.NotNil(t, response)
		assert.True(t, response.IsTokenNil())
	})

	t.Run("トークンの有効期限取得失敗", func(t *testing.T) {
		email := "test@mail.com"
		password := "password123"
		hashedPassword := "hashed_password"
		req := userUC.AuthenticateUserRequest{
			Email:    email,
			Password: password,
		}

		mockUserRepo.EXPECT().GetUserByEmail(context.Background(), email).Return(entity.NewUser(entity.UserParams{
			ID:         entity.UserID("user_id"),
			Name:       "John Doe",
			Email:      email,
			PasswdHash: hashedPassword,
			CreatedAt:  time.Now(),
			UpdatedAt:  nil,
		}), nil)

		mockHasher.EXPECT().ComparePassword(hashedPassword, password).Return(true, nil)
		mockTokenSvc.EXPECT().GenerateToken(gomock.Any()).Return("generated_token", nil)
		mockTokenSvc.EXPECT().GetExpireAt("generated_token").Return(0, errors.New("failed to get expiration time"))
		response, err := userUseCase.AuthenticateUser(context.Background(), req)
		assert.Error(t, err)
		assert.NotNil(t, response)
		assert.True(t, response.IsTokenNil())
		assert.Equal(t, response.GetExp(), 0)
	})
}
