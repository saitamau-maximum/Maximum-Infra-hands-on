package model

import (
	"time"

	"example.com/infrahandson/internal/domain/entity"
	"github.com/google/uuid"
)

type UserModel struct {
	ID           uuid.UUID  `db:"id"`
	Name         string     `db:"name"`
	Email        string     `db:"email"`
	PasswordHash string     `db:"password_hash"`
	ImagePath    *string     `db:"image_path"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    *time.Time `db:"updated_at"`
}

func (u *UserModel) ToEntity() *entity.User {
	return entity.NewUser(entity.UserParams{
		ID:         entity.UserID(u.ID.String()), // UUID -> UserID
		Name:       u.Name,
		Email:      u.Email,
		PasswdHash: u.PasswordHash,
		ImagePath:  u.ImagePath,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	})
}
