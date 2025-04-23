package adapter

import "example.com/webrtc-practice/internal/domain/entity"

type TokenServiceAdapter interface {
	GenerateToken(userPublicID entity.UserPublicID) (string, error)
	ValidateToken(token string) (string, error)
}
