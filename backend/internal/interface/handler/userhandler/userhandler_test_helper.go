package userhandler

import (
	"bytes"
	"io"
	"mime/multipart"

	"example.com/infrahandson/internal/infrastructure/validator"
	mock_adapter "example.com/infrahandson/test/mocks/interface/adapter"
	mock_factory "example.com/infrahandson/test/mocks/interface/factory"
	mock_usercase "example.com/infrahandson/test/mocks/usecase/usercase"
	"github.com/labstack/echo/v4"
	"go.uber.org/mock/gomock"
)

// mockDeps は UserHandler のテストで使用する依存関係モックをまとめた構造体です
type mockDeps struct {
	UserUseCase   mock_usercase.MockUserUseCaseInterface
	UserIDFactory mock_factory.MockUserIDFactory
	Logger        mock_adapter.MockLoggerAdapter
}

// NewTestUserHandler ( ハンドラ, モック依存関係, Echoインスタンス ) を生成する
func NewTestUserHandler(
	ctrl *gomock.Controller,
) (*UserHandler, mockDeps, *echo.Echo) {
	mockUserUseCase := mock_usercase.NewMockUserUseCaseInterface(ctrl)
	mockUserIDFactory := mock_factory.NewMockUserIDFactory(ctrl)
	mockLogger := mock_adapter.NewMockLoggerAdapter(ctrl)
	params := NewUserHandlerParams{
		UserUseCase:   mockUserUseCase,
		UserIDFactory: mockUserIDFactory,
		Logger:        mockLogger,
	}
	handler := NewUserHandler(params)

	mockDeps := mockDeps{
		UserUseCase:   *mockUserUseCase,
		UserIDFactory: *mockUserIDFactory,
		Logger:        *mockLogger,
	}

	e := echo.New()
	e.Validator = validator.NewEchoValidator()

	return handler, mockDeps, e
}

// NewMultipartFileHeader は， icon_test.go でダミーの画像リクエストを作成するためのヘルパー関数です．
func NewMultipartRequestWithFile(fieldName, filename string, content []byte, contentType string) (*bytes.Buffer, string, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	fw, err := w.CreateFormFile(fieldName, filename)
	if err != nil {
		return nil, "", err
	}
	if _, err := fw.Write(content); err != nil {
		return nil, "", err
	}

	if err := w.Close(); err != nil {
		return nil, "", err
	}

	return &b, w.FormDataContentType(), nil
}

// BadFile は、icon_test.go 用の不正なファイルを模倣するための構造体です
type BadFile struct{}

func (b *BadFile) Read(p []byte) (n int, err error)             { return 0, io.ErrUnexpectedEOF }
func (b *BadFile) Close() error                                 { return nil }
func (b *BadFile) Seek(offset int64, whence int) (int64, error) { return 0, nil }
