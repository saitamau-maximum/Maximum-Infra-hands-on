package adapter

type HasherAdapter interface {
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword, password string) error
}