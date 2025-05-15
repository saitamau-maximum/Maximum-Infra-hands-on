package service

import (
	"context"
	"io"

	"example.com/infrahandson/internal/domain/entity"
)

// iconデータやり取りのための構造体
type IconData struct {
	Reader   io.Reader
	Size     int64
	MimeType string
}

type IconStoreService interface {
	SaceIcon(ctx context.Context, iconData IconData, userID entity.UserID) error
	// 環境変数＋返り値pathにリダイレクトすることで画像を返す機構を想定
	GetIconPath(ctx context.Context, userID entity.UserID) (path string, err error)
}
