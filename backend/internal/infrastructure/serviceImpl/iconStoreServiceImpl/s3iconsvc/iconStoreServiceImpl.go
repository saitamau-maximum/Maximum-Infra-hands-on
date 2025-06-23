// s3を使用したアイコン保存の実装
package s3iconsvc

import (
	"bytes"
	"context"
	"errors"
	"image"
	"time"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/domain/service"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
	"github.com/chai2010/webp"
)

// S3またはMinIOを使用したアイコン保存の実装
// BaseURLは，S3とMinIOの共通ではない部分すべてを想定
// 例）
// S3の場合： https://examplebucket.s3.ap-northeast-1.amazonaws.com
// MinIOの場合： http://localhost:9000/my-bucket

type S3IconStoreImpl struct {
	baseURL string
	client *s3.Client
	bucket string
	prefix string
}

type NewS3IconStoreImplParams struct {
	BaseURL string
	Client *s3.Client
	Bucket string
	Prefix string
}

func (params *NewS3IconStoreImplParams) Validate() error {
	if params.BaseURL == "" {
		return errors.New("baseURL is required")
	}
	if params.Client == nil {
		return errors.New("client is required")
	}
	if params.Bucket == "" {
		return errors.New("bucket is required")
	}
	if params.Prefix == "" {
		return errors.New("key is required")
	}
	return nil
}

func NewS3IconStoreImpl(params NewS3IconStoreImplParams) service.IconStoreService {
	if err := params.Validate(); err != nil {
		panic(err)
	}

	return &S3IconStoreImpl{
		baseURL: params.BaseURL,
		client: params.Client,
		bucket: params.Bucket,
		prefix: params.Prefix,
	}
}

// アイコンアップロード処理。
// PutObject APIを使用して、io.Readerをバケットにアップロードする。
// 参考：https://docs.aws.amazon.com/ja_jp/code-library/latest/ug/go_2_s3_code_examples.html
// ↑アクション・PutObjectの例を参照
// https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/s3#PutObjectInput
// ↑PutObjectInputのフィールドを参照

func (s *S3IconStoreImpl) SaveIcon(ctx context.Context, iconData *service.IconData, userID entity.UserID) error {
	// サイズ検証
	if iconData.Size > service.GetMaxIconSize() {
		return errors.New("file size is too large")
	}

	// MIMEタイプ検証
	if iconData.MimeType != "image/jpeg" && iconData.MimeType != "image/png" && iconData.MimeType != "image/webp" {
		return errors.New("invalid mime type")
	}
	// アイコンの保存先のオブジェクトキーを生成
	objectKey := s.prefix + "/" + string(userID) + ".webp"

	//　アイコンをwebp形式に変換
	// イメージをデコード
	// 参考：https://pkg.go.dev/image#Decode
	img, _, err := image.Decode(bytes.NewReader(iconData.Icon))  // bytes.NewReaderで[]byteをReader化
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err = webp.Encode(&buf, img, &webp.Options{Lossless: true}); err != nil {
		return err
	}
	// アイコンを保存するためのPutObject APIを呼び出し
	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(s.bucket),
		Key:           aws.String(objectKey),
		Body:          bytes.NewReader(buf.Bytes()), // 変換したアイコンをio.Readerとしてアップロード
		ContentLength: aws.Int64(int64(buf.Len())),  // 変換後サイズを指定する
		ContentType:   aws.String("image/webp"),     // MIMEタイプを指定する
	})

	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "EntityTooLarge" {
			return errors.New("アイコンのサイズが大きすぎます")
		} else {
			return err
		}
	}

	// アイコン保存状態を確実にするためのNewObjectExistsWaiterを使用
	err = s3.NewObjectExistsWaiter(s.client).Wait(
		ctx, &s3.HeadObjectInput{Bucket: aws.String(s.bucket), Key: aws.String(objectKey)}, time.Minute)
	if err != nil {
		return err
	}

	return nil
}

func (s *S3IconStoreImpl) GetIconPath(ctx context.Context, userID entity.UserID) (string, error) {
	objectKey := s.prefix + "/" + string(userID) + ".webp"
	// リダイレクトできるように、URLを生成
	fullURL := s.baseURL + "/" + objectKey

	// オブジェクト存在確認
	_, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return "", err
	}
	// オブジェクトのURLを返す
	return fullURL, nil
}
