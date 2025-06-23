// ユーザーのアイコンを保存するロジックのインターフェース
// 具体実装は/interface/serviceImpl/iconStoreServiceImpl
package service

import (
	"context"

	"example.com/infrahandson/internal/domain/entity"
)

// MAX_ICON_SIZE はアイコンサイズでバリデーションするための定数
const MAX_ICON_SIZE = 5 * 1024 * 1024 // 5MB

// IconData はiconデータやり取りのための構造体（メソッドの入力）
type IconData struct {
	Icon   []byte
	Size     int64
	MimeType string
}

func NewIconData(reader []byte, size int64, mimeType string) *IconData {
	return &IconData{
		Icon:   reader,
		Size:     size,
		MimeType: mimeType,
	}
}

type IconStoreService interface {
	// SaveIcon はユーザーIDとアイコンデータ構造体を受け取り、保存する（パスも生成しておく）
	SaveIcon(ctx context.Context, iconData *IconData, userID entity.UserID) error

	// GetIconPath はUserIDからアイコンにアクセスするためのパス（URL）を返す
	// 環境変数＋返り値pathにリダイレクトすることで画像を返す機構を想定
	GetIconPath(ctx context.Context, userID entity.UserID) (path string, err error)
}

// GetMaxIconSize は定数を直接参照させないための関数
func GetMaxIconSize() int64 {
	return MAX_ICON_SIZE
}
