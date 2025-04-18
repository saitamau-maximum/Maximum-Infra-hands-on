package adapter

import "example.com/webrtc-practice/internal/domain/entity"

type TokenServiceAdapter interface {
	GenerateToken(userID entity.UserID) (string, error)
	ValidateToken(token string) (int, error)
}
