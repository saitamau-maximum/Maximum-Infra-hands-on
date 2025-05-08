package repository

import (
	"context"

	"example.com/infrahandson/internal/domain/entity"
)

type UserRepository interface {
	SaveUser(context.Context, *entity.User) (*entity.User, error)
	GetUserByID(ctx context.Context, id entity.UserID) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
}
