package localiconstoreimpl

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/domain/service"
	"github.com/pkg/errors"
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
	// ファイル名はユーザーのUUIDにする
	fileName := string(userID) + ".webp"
	path := filepath.Join(l.dirPath, fileName)

	// 書き込み先のファイルを開く
	dst, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "failed to create file")
	}
	defer dst.Close()

	// サイズ検証
	if iconData.Size > service.GetMaxIconSize() {
		return errors.New("file size is too large")
	}

	// MIMEタイプ検証
	if iconData.MimeType != "image/jpeg" && iconData.MimeType != "image/png" && iconData.MimeType != "image/webp" {
		return errors.New("invalid mime type")
	}

	// 書き込み
	if _, err := io.Copy(dst, iconData.Reader); err != nil {
		return errors.Wrap(err, "failed to copy file")
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
	
	return path, nil
}
