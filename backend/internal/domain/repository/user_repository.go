package repository

import (
	"example.com/webrtc-practice/internal/domain/entity"
)

type UserRepository interface {
	CreateUser(entity.User) (*entity.User, error)
	GetAllUsers() ([]*entity.User, error)
	GetUserByID(id int) (*entity.User, error)
	GetUserByEmail(email string) (*entity.User, error)
	UpdateUser(user *entity.User) (*entity.User, error)
	DeleteUser(id int) error
}
