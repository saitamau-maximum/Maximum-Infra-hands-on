package service

import (
	"context"

	"example.com/infrahandson/internal/domain/entity"
)

type ImageStoreService interface {
	SaveImage(ctx context.Context, imageData []byte, userID entity.UserID) (path string, err error)
	// Update the image without changing its path
	UpdateImage(ctx context.Context, imageData []byte, userID entity.UserID) error
}
