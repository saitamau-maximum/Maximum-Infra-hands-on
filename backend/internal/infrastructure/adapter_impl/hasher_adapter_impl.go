package adapter_impl

import (
	"errors"

	"example.com/webrtc-practice/internal/interface/adapter"
	"golang.org/x/crypto/bcrypt"
)

type HasherAdapterImpl struct {
	cost int
}

type NewHasherAddapterParams struct {
	Cost int
}

func (p *NewHasherAddapterParams) Validate() error {
	// cost must be between 1 and 31
	if p.Cost <= 0 || p.Cost > 31 {
		return errors.New("cost must be greater than 0")
	}
	return nil
}

func NewHasherAdapter(params NewHasherAddapterParams) adapter.HasherAdapter {
	if err := params.Validate(); err != nil {
		panic(err)
	}

	return &HasherAdapterImpl{
		cost: params.Cost,
	}
}

func (h *HasherAdapterImpl) HashPassword(password string) (string, error) {
	// Hash the password using bcrypt with the specified cost
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (h *HasherAdapterImpl) ComparePassword(hashedPassword, password string) (bool, error) {
	// Compare the hashed password with the provided password
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil // Passwords do not match
		}
		return false, err // Some other error occurred
	}
	return true, nil // Passwords match
}
