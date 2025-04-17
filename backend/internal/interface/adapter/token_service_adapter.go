package adapter

type TokenServiceAdapter interface {
	GenerateToken(userID int) (string, error)
	ValidateToken(token string) (int, error)
}