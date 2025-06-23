// ユーザーの永続化に関わるインターフェース
// 具体実装は/infrastructure/repositoryImpl
package repository

import (
	"context"

	"example.com/infrahandson/internal/domain/entity"
)

type UserRepository interface {
	// SaveUserはユーザー情報を保存し、保存後のユーザーを返します。
	SaveUser(ctx context.Context, user *entity.User) (*entity.User, error)

	// GetUserByIDは指定したユーザーIDに対応するユーザー情報を取得します。
	GetUserByID(ctx context.Context, id entity.UserID) (*entity.User, error)

	// GetUserByEmailは指定したメールアドレスに対応するユーザー情報を取得します。
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
}
