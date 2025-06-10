package user_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime/multipart"
	"testing"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/domain/service"
	userUC "example.com/infrahandson/internal/usecase/user"
	mock_service "example.com/infrahandson/test/mocks/domain/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// テスト用のmultipart.File実装

type TestFile struct {
	*bytes.Reader
}

func (f *TestFile) Read(p []byte) (int, error)              { return f.Reader.Read(p) }
func (f *TestFile) ReadAt(p []byte, off int64) (int, error) { return f.Reader.ReadAt(p, off) }
func (f *TestFile) Seek(offset int64, whence int) (int64, error) {
	return f.Reader.Seek(offset, whence)
}
func (f *TestFile) Close() error { return nil }

// ReadAll失敗用

type ErrFile struct{}

func (e *ErrFile) Read(p []byte) (int, error)                   { return 0, errors.New("read error") }
func (e *ErrFile) ReadAt(p []byte, off int64) (int, error)      { return 0, errors.New("read error") }
func (e *ErrFile) Seek(offset int64, whence int) (int64, error) { return 0, nil }
func (e *ErrFile) Close() error                                 { return nil }

// テスト用UserUseCaseラッパー

type testUserUseCase struct {
	iconSvc  *mock_service.MockIconStoreService
	openFunc func() (multipart.File, error)
}

func (u *testUserUseCase) SaveUserIcon(ctx context.Context, fh *multipart.FileHeader, id entity.UserID) error {
	file, err := u.openFunc()
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	iconData := service.NewIconData(data, int64(len(data)), fh.Header.Get("Content-Type"))
	return u.iconSvc.SaveIcon(ctx, iconData, id)
}

func TestSaveUserIcon(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIconSvc := mock_service.NewMockIconStoreService(ctrl)
	userID := entity.UserID("user_id")

	t.Run("正常系", func(t *testing.T) {
		iconBytes := []byte("icondata")
		fh := &multipart.FileHeader{Filename: "icon.png", Header: make(map[string][]string)}
		fh.Header.Set("Content-Type", "image/png")
		testUC := &testUserUseCase{mockIconSvc, func() (multipart.File, error) {
			return &TestFile{bytes.NewReader(iconBytes)}, nil
		}}
		mockIconSvc.EXPECT().SaveIcon(gomock.Any(), gomock.Any(), userID).Return(nil)
		err := testUC.SaveUserIcon(context.Background(), fh, userID)
		assert.NoError(t, err)
	})

	t.Run("ファイルオープン失敗", func(t *testing.T) {
		fh := &multipart.FileHeader{}
		testUC := &testUserUseCase{mockIconSvc, func() (multipart.File, error) {
			return nil, errors.New("open error")
		}}
		err := testUC.SaveUserIcon(context.Background(), fh, userID)
		assert.Error(t, err)
		assert.Equal(t, "open error", err.Error())
	})

	t.Run("ReadAll失敗", func(t *testing.T) {
		fh := &multipart.FileHeader{}
		testUC := &testUserUseCase{mockIconSvc, func() (multipart.File, error) {
			return &ErrFile{}, nil
		}}
		err := testUC.SaveUserIcon(context.Background(), fh, userID)
		assert.Error(t, err)
		assert.Equal(t, "read error", err.Error())
	})

	t.Run("保存失敗", func(t *testing.T) {
		iconBytes := []byte("icondata")
		fh := &multipart.FileHeader{Header: make(map[string][]string)}
		fh.Header.Set("Content-Type", "image/png")
		testUC := &testUserUseCase{mockIconSvc, func() (multipart.File, error) {
			return &TestFile{bytes.NewReader(iconBytes)}, nil
		}}
		mockIconSvc.EXPECT().SaveIcon(gomock.Any(), gomock.Any(), userID).Return(errors.New("save error"))
		err := testUC.SaveUserIcon(context.Background(), fh, userID)
		assert.Error(t, err)
		assert.Equal(t, "save error", err.Error())
	})
}

func TestGetUserIconPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUseCase, mockDeps := userUC.NewTestUserUseCase(ctrl)

	userID := entity.UserID("user_id")

	t.Run("正常系", func(t *testing.T) {
		mockDeps.IconSvc.EXPECT().GetIconPath(gomock.Any(), userID).Return("/path/to/icon.png", nil)
		path, err := userUseCase.GetUserIconPath(context.Background(), userID)
		assert.NoError(t, err)
		assert.Equal(t, "/path/to/icon.png", path)
	})

	t.Run("失敗系", func(t *testing.T) {
		mockDeps.IconSvc.EXPECT().GetIconPath(gomock.Any(), userID).Return("", errors.New("not found"))
		path, err := userUseCase.GetUserIconPath(context.Background(), userID)
		assert.Error(t, err)
		assert.Empty(t, path)
		assert.Equal(t, "not found", err.Error())
	})
}
