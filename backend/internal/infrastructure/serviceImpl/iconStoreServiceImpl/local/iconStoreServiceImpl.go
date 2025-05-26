package localiconstoreimpl

import (
	"bytes"
	"context"
	"errors"
	"image"
	"os"
	"path/filepath"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/domain/service"

	"github.com/chai2010/webp"

	_ "image/jpeg"
	_ "image/png"
)

type localiconstoreimpl struct {
	dirPath string
}

type NewLocalIconStoreImplParams struct {
	DirPath string
}

func (p *NewLocalIconStoreImplParams) Validate() error {
	if p.DirPath == "" {
		return errors.New("dir is required")
	}
	return nil
}

func NewLocalIconStoreImpl(p *NewLocalIconStoreImplParams) service.IconStoreService {
	if err := p.Validate(); err != nil {
		panic(err)
	}

	// 存在しないなら、ディレクトリを作成する
	if err := os.MkdirAll(p.DirPath, 0755); err != nil {
		panic(err)
	}

	return &localiconstoreimpl{
		dirPath: p.DirPath,
	}
}

func (l *localiconstoreimpl) SaveIcon(ctx context.Context, iconData *service.IconData, userID entity.UserID) error {
	// サイズ検証
	if iconData.Size > service.GetMaxIconSize() {
		return errors.New("file size is too large")
	}

	// MIMEタイプ検証
	if iconData.MimeType != "image/jpeg" && iconData.MimeType != "image/png" && iconData.MimeType != "image/webp" {
		return errors.New("invalid mime type")
	}

	// ファイル名はユーザーのUUIDにする
	fileName := string(userID) + ".webp"
	path := filepath.Join(l.dirPath, fileName)

	// イメージをデコード（[]byte を io.Reader に変換）
	img, _, err := image.Decode(bytes.NewReader(iconData.Icon))
	if err != nil {
		return errors.New("failed to decode image: " + err.Error())
	}
	// 書き込み先のファイルを開く
	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()

	// 書き込み先のファイルにエンコード
	// 参考: https://github.com/chai2010/webp?tab=readme-ov-file#example
	if err := webp.Encode(dst, img, &webp.Options{Quality: 75}); err != nil {
		return err
	}
	return nil
}

func (l *localiconstoreimpl) GetIconPath(ctx context.Context, userID entity.UserID) (string, error) {
	// ファイル名はユーザーのUUIDにする
	fileName := string(userID) + ".webp"
	path := filepath.Join(l.dirPath, fileName)
	// ファイルが存在しない場合はエラーを返す
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", errors.New("file not found")
	}
	// リダイレクトを絶対パスにする
	path = "/" + path
	return path, nil
}
