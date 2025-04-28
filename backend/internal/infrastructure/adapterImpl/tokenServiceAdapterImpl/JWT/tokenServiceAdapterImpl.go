package tokenadapterimpl

import (
	"errors"
	"time"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/interface/adapter"
	"github.com/golang-jwt/jwt/v5"
)

type TokenServiceAdapterImpl struct {
	secretKey     string
	expireMinutes int
}

type NewTokenServiceAdapterParams struct {
	SecretKey     string
	ExpireMinutes int
}

func (p *NewTokenServiceAdapterParams) Validate() error {
	if p.SecretKey == "" {
		return errors.New("secret key must not be empty")
	}
	if p.ExpireMinutes <= 0 {
		return errors.New("expireMinutes must be greater than 0")
	}
	return nil
}

func NewTokenServiceAdapter(params NewTokenServiceAdapterParams) adapter.TokenServiceAdapter {
	if err := params.Validate(); err != nil {
		panic(err)
	}
	return &TokenServiceAdapterImpl{
		secretKey:     params.SecretKey,
		expireMinutes: params.ExpireMinutes,
	}
}

func (s *TokenServiceAdapterImpl) GenerateToken(userPublicID entity.UserPublicID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": string(userPublicID),
		"exp":     time.Now().Add(time.Duration(s.expireMinutes) * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

func (s *TokenServiceAdapterImpl) ValidateToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.secretKey), nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("user_id not found in token or is not a string")
	}

	return userIDStr, nil
}

func (s *TokenServiceAdapterImpl) GetExpireAt(tokenStr string) (int, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		exp, ok := claims["exp"].(float64)
		if !ok {
			return 0, errors.New("expiration time not found in token")
		}
		return int(exp), nil
	}

	return 0, errors.New("invalid token")
}
