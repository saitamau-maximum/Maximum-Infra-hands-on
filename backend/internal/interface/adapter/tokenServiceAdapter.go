// トークン関連の機能アダプター
// 具体実装は/infrastructure/adapterImpl/tokenServiceAdapterImpl
package adapter

import "example.com/infrahandson/internal/domain/entity"

type TokenServiceAdapter interface {
	// GenerateToken はユーザーIDからトークンを生成します。
	GenerateToken(userID entity.UserID) (string, error)

	// ValidateToken はトークンを検証し、ユーザーIDを返します。
	ValidateToken(token string) (userID string, err error)

	// GetExpireAt はトークンの有効期限を取得します。
	GetExpireAt(token string) (int, error)
}
