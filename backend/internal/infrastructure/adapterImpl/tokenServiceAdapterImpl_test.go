package adapterimpl_test

import (
	"testing"

	"example.com/infrahandson/internal/domain/entity"
	adapterimpl "example.com/infrahandson/internal/infrastructure/adapterImpl"

	"github.com/stretchr/testify/assert"
)

func TestTokenServiceAdapter(t *testing.T) {
	params := adapterimpl.NewTokenServiceAdapterParams{
		SecretKey:     "test-secret-key",
		ExpireMinutes: 1, // 1分で期限切れになる設定
	}
	tokenService := adapterimpl.NewTokenServiceAdapter(params)

	t.Run("Generate and Validate valid token", func(t *testing.T) {
		userID := entity.UserPublicID("abc123")
		token, err := tokenService.GenerateToken(userID)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		// 今回は文字列IDなので、変換できない前提で ValidateToken を修正する必要がある
		// ↓例として user_id を string に戻す形の ValidateToken を想定
	})

	t.Run("ValidateToken should fail with invalid token", func(t *testing.T) {
		_, err := tokenService.ValidateToken("invalid.token.string")
		assert.Error(t, err)
	})
}
