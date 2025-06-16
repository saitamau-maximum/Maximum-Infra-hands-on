package usercase

import (
	"context"
	"errors"
)

// AuthenticateUserRequest構造体: 認証リクエスト
type AuthenticateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthenticateUserResponse構造体: 認証レスポンス
type AuthenticateUserResponse struct {
	token *string
	exp   *int
}

// IsTokenNil: トークンがnilかどうかを判定
func (res *AuthenticateUserResponse) IsTokenNil() bool {
	return res.token == nil
}

// GetToken: トークンを取得（nilの場合は空文字を返す）
func (res *AuthenticateUserResponse) GetToken() string {
	if res.token == nil {
		return ""
	}
	return *res.token
}

func (res *AuthenticateUserResponse) GetExp() int {
	if res.exp == nil {
		return 0
	}
	return *res.exp
}

// 外部でのテストのためのセッター
func (res *AuthenticateUserResponse) SetToken(token string) {
	res.token = &token
}

func (res *AuthenticateUserResponse) SetExp(exp int) {
	res.exp = &exp
}

// AuthenticateUser ユーザー認証
func (u *UserUseCase) AuthenticateUser(ctx context.Context, req AuthenticateUserRequest) (AuthenticateUserResponse, error) {
	user, err := u.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return AuthenticateUserResponse{token: nil}, err
	}

	ok, err := u.hasher.ComparePassword(user.GetPasswdHash(), req.Password)
	if err != nil {
		return AuthenticateUserResponse{token: nil}, err
	}

	if !ok {
		return AuthenticateUserResponse{token: nil}, errors.New("password mismatch")
	}

	res, err := u.tokenSvc.GenerateToken(user.GetID())
	if err != nil {
		return AuthenticateUserResponse{token: nil}, err
	}

	exp, err := u.tokenSvc.GetExpireAt(res)
	if err != nil {
		return AuthenticateUserResponse{token: nil}, err
	}

	return AuthenticateUserResponse{
		token: &res,
		exp:   &exp,
	}, nil
}
