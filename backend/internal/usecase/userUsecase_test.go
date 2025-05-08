package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/usecase"
	mock_repository "example.com/infrahandson/test/mocks/domain/repository"
	mock_adapter "example.com/infrahandson/test/mocks/interface/adapter"
	mock_factory "example.com/infrahandson/test/mocks/interface/factory"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSignUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// == Mock dependencies ==
	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockHasher := mock_adapter.NewMockHasherAdapter(ctrl)
	mockTokenSvc := mock_adapter.NewMockTokenServiceAdapter(ctrl)
	mockUserIDFactory := mock_factory.NewMockUserIDFactory(ctrl)

	params := usecase.NewUserUseCaseParams{
		UserRepo:      mockUserRepo,
		Hasher:        mockHasher,
		TokenSvc:      mockTokenSvc,
		UserIDFactory: mockUserIDFactory,
	}

	userUseCase := usecase.NewUserUseCase(params)

	t.Run("正常系", func(t *testing.T) {
		signUpRequest := usecase.SignUpRequest{
			Name:     "John Doe",
			Email:    "test@mail.com",
			Password: "password123",
		}
		hashedPassword := "hashed_password"
		userID := entity.UserID("public_user_id")

		mockHasher.EXPECT().HashPassword(signUpRequest.Password).Return(hashedPassword, nil)
		mockUserIDFactory.EXPECT().NewUserID().Return(userID, nil)
		mockUserRepo.EXPECT().
			SaveUser(context.Background(), gomock.AssignableToTypeOf(&entity.User{})).
			DoAndReturn(func(u *entity.User) (*entity.User, error) {
				return u, nil
			})

		response, err := userUseCase.SignUp(context.Background(), signUpRequest)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, signUpRequest.Name, response.User.GetName())
		assert.Equal(t, signUpRequest.Email, response.User.GetEmail())
		assert.Equal(t, hashedPassword, response.User.GetPasswdHash())
		assert.Equal(t, userID, response.User.GetID())
	})

	t.Run("ハッシュ化失敗時", func(t *testing.T) {
		signUpRequest := usecase.SignUpRequest{
			Name:     "John Doe",
			Email:    "test@mail.com",
			Password: "password123",
		}
		expectedErr := errors.New("failed to hash password")

		mockHasher.EXPECT().HashPassword(signUpRequest.Password).Return("", expectedErr)

		response, err := userUseCase.SignUp(context.Background(), signUpRequest)

		assert.Error(t, err)
		assert.Nil(t, response.User)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("UserID生成失敗時", func(t *testing.T) {
		signUpRequest := usecase.SignUpRequest{
			Name:     "John Doe",
			Email:    "test@mail.com",
			Password: "password123",
		}
		hashedPassword := "hashed_password"
		expectedErr := errors.New("failed to generate user ID")

		mockHasher.EXPECT().HashPassword(signUpRequest.Password).Return(hashedPassword, nil)
		mockUserIDFactory.EXPECT().NewUserID().Return(entity.UserID(""), expectedErr)

		response, err := userUseCase.SignUp(context.Background(), signUpRequest)

		assert.Error(t, err)
		assert.Nil(t, response.User)
		assert.Equal(t, expectedErr, err)
	})
}

func TestAuthenticateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// == Mock dependencies ==
	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockHasher := mock_adapter.NewMockHasherAdapter(ctrl)
	mockTokenSvc := mock_adapter.NewMockTokenServiceAdapter(ctrl)
	mockUserIDFactory := mock_factory.NewMockUserIDFactory(ctrl)

	params := usecase.NewUserUseCaseParams{
		UserRepo:      mockUserRepo,
		Hasher:        mockHasher,
		TokenSvc:      mockTokenSvc,
		UserIDFactory: mockUserIDFactory,
	}

	userUseCase := usecase.NewUserUseCase(params)

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

		req := usecase.AuthenticateUserRequest{
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

	t.Run("ユーザーが存在しない場合", func(t *testing.T) {
		email := "notfound@email.com"
		password := "password123"

		req := usecase.AuthenticateUserRequest{
			Email:    email,
			Password: password,
		}

		mockUserRepo.EXPECT().GetUserByEmail(context.Background(), email).Return(nil, errors.New("user not found"))

		response, err := userUseCase.AuthenticateUser(context.Background(), req)
		assert.Error(t, err)
		assert.Equal(t, "", response.GetToken())
	})

	t.Run("パスワードが一致しない場合", func(t *testing.T) {
		email := "test@email.com"
		password := "wrongpassword"
		hashedPassword := "correct_hash"

		user := entity.NewUser(entity.UserParams{
			ID:         entity.UserID("user_id"),
			Name:       "John Doe",
			Email:      email,
			PasswdHash: hashedPassword,
			CreatedAt:  time.Now(),
			UpdatedAt:  nil,
		})

		req := usecase.AuthenticateUserRequest{
			Email:    email,
			Password: password,
		}

		mockUserRepo.EXPECT().GetUserByEmail(context.Background(), email).Return(user, nil)
		mockHasher.EXPECT().ComparePassword(hashedPassword, password).Return(false, nil)

		response, err := userUseCase.AuthenticateUser(context.Background(), req)
		assert.Error(t, err)
		assert.Equal(t, "", response.GetToken())
	})
}
