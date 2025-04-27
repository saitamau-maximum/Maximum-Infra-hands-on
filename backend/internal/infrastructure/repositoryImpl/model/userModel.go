package model

import (
	"time"

	"example.com/infrahandson/internal/domain/entity"
)

type UserModel struct {
	ID           int        `db:"id"`
	PublicID     string     `db:"public_id"`
	Name         string     `db:"name"`
	Email        string     `db:"email"`
	PasswordHash string     `db:"password_hash"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    *time.Time `db:"updated_at"`
}

func (u *UserModel) ToEntity() *entity.User {
	return entity.NewUser(entity.UserParams{
		ID:         entity.UserID(u.ID),
		PublicID:   entity.UserPublicID(u.PublicID),
		Name:       u.Name,
		Email:      u.Email,
		PasswdHash: u.PasswordHash,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	})
}
