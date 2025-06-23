package usercase

import (
	"context"
	"mime/multipart"

	"example.com/infrahandson/internal/domain/entity"
)

// UserUseCaseInterface: ユーザーに関するユースケースを管理するインターフェース
type UserUseCaseInterface interface {
	// SignUp: ユーザー登録を行う
	SignUp(ctx context.Context, req SignUpRequest) (SignUpResponse, error)

	// AuthenticateUser: ユーザー認証を行う
	AuthenticateUser(ctx context.Context, req AuthenticateUserRequest) (AuthenticateUserResponse, error)

	// GetUserByID: ユーザーIDからユーザー情報を取得する
	GetUserByID(ctx context.Context, id entity.UserID) (*entity.User, error)

	// SaveUserIcon: ユーザーのアイコンを保存する
	SaveUserIcon(ctx context.Context, fh *multipart.FileHeader, id entity.UserID) error

	// SaveUserIcon: ユーザーのアイコンを保存する
	GetUserIconPath(ctx context.Context, id entity.UserID) (path string, err error)
}
