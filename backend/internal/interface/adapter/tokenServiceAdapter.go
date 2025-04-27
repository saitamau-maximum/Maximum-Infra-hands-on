package adapter

import "example.com/infrahandson/internal/domain/entity"

type TokenServiceAdapter interface {
	GenerateToken(userPublicID entity.UserPublicID) (string, error)
	ValidateToken(token string) (string, error)
	GetExpireAt(token string) (int, error)
}
