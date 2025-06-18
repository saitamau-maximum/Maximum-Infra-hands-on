package bcrypt_test

import (
	"testing"

	bcryptadapterimpl "example.com/infrahandson/internal/infrastructure/adapterImpl/hasherAdapterImpl/bcrypt"
	"github.com/stretchr/testify/assert"
)

func TestHasherAdapterImpl_HashAndComparePassword(t *testing.T) {
	adapter := bcryptadapterimpl.NewHasherAdapter(bcryptadapterimpl.NewHasherAddapterParams{Cost: 10})

	password := "securePassword123"

	t.Run("正常系: ハッシュ化して比較一致", func(t *testing.T) {
		hashed, err := adapter.HashPassword(password)
		assert.NoError(t, err)
		assert.NotEmpty(t, hashed)

		match, err := adapter.ComparePassword(hashed, password)
		assert.NoError(t, err)
		assert.True(t, match)
	})

	t.Run("異常系: 間違ったパスワードで照合", func(t *testing.T) {
		hashed, err := adapter.HashPassword(password)
		assert.NoError(t, err)

		match, err := adapter.ComparePassword(hashed, "wrongPassword")
		assert.NoError(t, err)
		assert.False(t, match)
	})
}

func TestNewHasherAdapter_InvalidCost(t *testing.T) {
	t.Run("costが0以下だとpanic", func(t *testing.T) {
		assert.Panics(t, func() {
			_ = bcryptadapterimpl.NewHasherAdapter(bcryptadapterimpl.NewHasherAddapterParams{Cost: 0})
		})
	})

	t.Run("costが31を超えるとpanic", func(t *testing.T) {
		assert.Panics(t, func() {
			_ = bcryptadapterimpl.NewHasherAdapter(bcryptadapterimpl.NewHasherAddapterParams{Cost: 32})
		})
	})
}
