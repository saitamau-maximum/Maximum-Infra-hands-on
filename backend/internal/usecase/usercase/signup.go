package usercase

import (
	"context"
	"time"

	"example.com/infrahandson/internal/domain/entity"
)

// SignUpRequest構造体: サインアップリクエスト
type SignUpRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SignUpResponse構造体: サインアップレスポンス
type SignUpResponse struct {
	User *entity.User
}

// SignUp ユーザー登録
func (u *UserUseCase) SignUp(ctx context.Context, req SignUpRequest) (SignUpResponse, error) {
	hashedPassword, err := u.hasher.HashPassword(req.Password)
	if err != nil {
		return SignUpResponse{nil}, err
	}

	id, err := u.userIDFactory.NewUserID()
	if err != nil {
		return SignUpResponse{nil}, err
	}

	userParams := entity.UserParams{
		ID:         id,
		Name:       req.Name,
		Email:      req.Email,
		PasswdHash: hashedPassword,
		CreatedAt:  time.Now(),
		UpdatedAt:  nil,
	}

	user := entity.NewUser(userParams)

	res, err := u.userRepo.SaveUser(ctx, user)
	if err != nil {
		return SignUpResponse{nil}, err
	}

	return SignUpResponse{User: res}, nil
}
