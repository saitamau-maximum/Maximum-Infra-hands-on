package entity

import (
	"time"
)
type User struct {
	id         UserID
	name       string
	email      string
	passwdhash string
	createdAt  time.Time
	updatedAt  *time.Time
}

func NewUser(id UserID, name, email, passwdhash string, createdAt time.Time, updatedAt *time.Time) *User {
	return &User{
		id:         id,
		name:       name,
		email:      email,
		passwdhash: passwdhash,
		createdAt:  createdAt,
		updatedAt:  updatedAt,
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
