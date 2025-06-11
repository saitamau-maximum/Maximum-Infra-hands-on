package usercase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/usecase/usercase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// 1.正常系
// 2.GetUserByID失敗

func TestGetUserByID(t *testing.T) {
	// Initialize the mock objects
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUseCase, mockDeps := usercase.NewTestUserUseCase(ctrl)

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
		mockDeps.UserRepo.EXPECT().GetUserByID(gomock.Any(), entity.UserID("user_id")).Return(user, nil)
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
		mockDeps.UserRepo.EXPECT().GetUserByID(gomock.Any(), entity.UserID("invalid_id")).Return(nil, errors.New("user not found"))
		resp, err := userUseCase.GetUserByID(context.Background(), entity.UserID("invalid_id"))
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}
