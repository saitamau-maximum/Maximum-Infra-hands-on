package model

type UserModel struct {
	ID           int     `db:"id"`
	PublicID     string  `db:"public_id"`
	Name         string  `db:"name"`
	Email        string  `db:"email"`
	PasswordHash string  `db:"password_hash"`
	CreatedAt    string  `db:"created_at"`
	UpdatedAt    *string `db:"updated_at"`
}
