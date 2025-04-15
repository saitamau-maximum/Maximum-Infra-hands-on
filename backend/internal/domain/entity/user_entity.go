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



