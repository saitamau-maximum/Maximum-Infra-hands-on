package service

import (
	"context"
	"io"

	"example.com/infrahandson/internal/domain/entity"
)

const MAX_ICON_SIZE = 5 * 1024 * 1024 // 5MB

// iconデータやり取りのための構造体
type IconData struct {
	Reader   io.Reader
	Size     int64
	MimeType string
}

func NewIconData(reader io.Reader, size int64, mimeType string) *IconData {
	return &IconData{
		Reader:   reader,
		Size:     size,
		MimeType: mimeType,
	}
}

type IconStoreService interface {
	SaveIcon(ctx context.Context, iconData *IconData, userID entity.UserID) error
	// 環境変数＋返り値pathにリダイレクトすることで画像を返す機構を想定
	GetIconPath(ctx context.Context, userID entity.UserID) (path string, err error)
}

func GetMaxIconSize() int64 {
	return MAX_ICON_SIZE
}
