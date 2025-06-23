package usercase

import (
	"context"

	"example.com/infrahandson/internal/domain/entity"
)

func (u *UserUseCase) GetUserByID(ctx context.Context, id entity.UserID) (*entity.User, error) {
	user, err := u.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
