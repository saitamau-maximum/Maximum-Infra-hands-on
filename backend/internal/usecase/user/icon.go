package user

import (
	"context"
	"io"
	"mime/multipart"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/domain/service"
)

func (u *UserUseCase) SaveUserIcon(ctx context.Context, fh *multipart.FileHeader, id entity.UserID) error {
	// ファイルを開く
	file, err := fh.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	// []byte に読み込む
	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	// IconData を作成
	iconData := service.NewIconData(data, int64(len(data)), fh.Header.Get("Content-Type"))

	// 保存処理を呼ぶ
	return u.iconSvc.SaveIcon(ctx, iconData, id)
}

func (u *UserUseCase) GetUserIconPath(ctx context.Context, id entity.UserID) (path string, err error) {
	path, err = u.iconSvc.GetIconPath(ctx, id)
	if err != nil {
		return "", err
	}
	return path, nil
}
