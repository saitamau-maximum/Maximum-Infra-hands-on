package user_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"example.com/infrahandson/internal/domain/entity"
	userUC "example.com/infrahandson/internal/usecase/user"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// TestSignUp ユーザー登録のテスト
// 1. 正常系: ユーザー登録が成功することを確認
// 2. ハッシュ化失敗時: パスワードのハッシュ化に失敗する場合の挙動を確認
// 3. UserID生成失敗時: ユーザーIDの生成に失敗する場合の挙動を確認
// 4. ユーザー保存失敗時: ユーザーの保存に失敗する場合の挙動を確認

func TestSignUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUseCase, mockDeps := userUC.NewTestUserUseCase(ctrl)

	t.Run("正常系", func(t *testing.T) {
		signUpRequest := userUC.SignUpRequest{
			Name:     "John Doe",
			Email:    "test@mail.com",
			Password: "password123",
		}
		hashedPassword := "hashed_password"
		userID := entity.UserID("public_user_id")
		expectUser := entity.NewUser(entity.UserParams{
			ID:         userID,
			Name:       signUpRequest.Name,
			Email:      signUpRequest.Email,
			PasswdHash: hashedPassword,
			CreatedAt:  time.Now(),
			UpdatedAt:  nil,
		})

		mockDeps.Hasher.EXPECT().HashPassword(signUpRequest.Password).Return(hashedPassword, nil)
		mockDeps.UserIDFactory.EXPECT().NewUserID().Return(userID, nil)
		mockDeps.UserRepo.EXPECT().SaveUser(context.Background(), gomock.Any()).Return(expectUser, nil)

		response, err := userUseCase.SignUp(context.Background(), signUpRequest)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, signUpRequest.Name, response.User.GetName())
		assert.Equal(t, signUpRequest.Email, response.User.GetEmail())
		assert.Equal(t, hashedPassword, response.User.GetPasswdHash())
		assert.Equal(t, userID, response.User.GetID())
	})

	t.Run("ハッシュ化失敗時", func(t *testing.T) {
		signUpRequest := userUC.SignUpRequest{
			Name:     "John Doe",
			Email:    "test@mail.com",
			Password: "password123",
		}
		expectedErr := errors.New("failed to hash password")

		mockDeps.Hasher.EXPECT().HashPassword(signUpRequest.Password).Return("", expectedErr)

		response, err := userUseCase.SignUp(context.Background(), signUpRequest)

		assert.Error(t, err)
		assert.Nil(t, response.User)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("UserID生成失敗時", func(t *testing.T) {
		signUpRequest := userUC.SignUpRequest{
			Name:     "John Doe",
			Email:    "test@mail.com",
			Password: "password123",
		}
		hashedPassword := "hashed_password"
		expectedErr := errors.New("failed to generate user ID")

		mockDeps.Hasher.EXPECT().HashPassword(signUpRequest.Password).Return(hashedPassword, nil)
		mockDeps.UserIDFactory.EXPECT().NewUserID().Return(entity.UserID(""), expectedErr)

		response, err := userUseCase.SignUp(context.Background(), signUpRequest)

		assert.Error(t, err)
		assert.Nil(t, response.User)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("ユーザー保存失敗時", func(t *testing.T) {
		signUpRequest := userUC.SignUpRequest{
			Name:     "John Doe",
			Email:    "test@mail.com",
			Password: "password123",
		}
		hashedPassword := "hashed_password"
		userID := entity.UserID("public_user_id")
		expectedErr := errors.New("failed to save user")
		mockDeps.Hasher.EXPECT().HashPassword(signUpRequest.Password).Return(hashedPassword, nil)
		mockDeps.UserIDFactory.EXPECT().NewUserID().Return(userID, nil)
		mockDeps.UserRepo.EXPECT().SaveUser(context.Background(), gomock.Any()).Return(nil, expectedErr)
		response, err := userUseCase.SignUp(context.Background(), signUpRequest)
		assert.Error(t, err)
		assert.Nil(t, response.User)
		assert.Equal(t, expectedErr, err)
	})
}
