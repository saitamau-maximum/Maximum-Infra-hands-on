package entity

import (
	"time"
)

type User struct {
	id         UserID
	name       string
	email      string
	passwdhash string
	imagePath  *string
	createdAt  time.Time
	updatedAt  *time.Time
}

type UserParams struct {
	ID         UserID
	Name       string
	Email      string
	PasswdHash string
	ImagePath  *string
	CreatedAt  time.Time
	UpdatedAt  *time.Time
}

func NewUser(p UserParams) *User {
	return &User{
		id:         p.ID,
		name:       p.Name,
		email:      p.Email,
		passwdhash: p.PasswdHash,
		imagePath:  p.ImagePath,
		createdAt:  p.CreatedAt,
		updatedAt:  p.UpdatedAt,
	}
}

func (u User) GetID() UserID {
	return u.id
}

func (u User) GetName() string {
	return u.name
}

func (u User) GetEmail() string {
	return u.email
}

func (u User) GetPasswdHash() string {
	return u.passwdhash
}

func (u User) GetCreatedAt() time.Time {
	return u.createdAt
}

func (u User) GetImagePath() (imagePath string, ok bool) {
	if u.imagePath == nil {
		return "", false
	}
	return *u.imagePath, true
}

func (u *User) SetImagePath(imagePath string) {
	u.imagePath = &imagePath
}