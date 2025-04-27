package repository

import "example.com/infrahandson/internal/domain/entity"

type UserRepository interface {
	SaveUser(*entity.User) (*entity.User, error)
	GetUserByID(id entity.UserID) (*entity.User, error)
	GetUserByEmail(email string) (*entity.User, error)
	GetIDByPublicID(publicID entity.UserPublicID) (entity.UserID, error)
	GetPublicIDByID(id entity.UserID) (entity.UserPublicID, error)
}
