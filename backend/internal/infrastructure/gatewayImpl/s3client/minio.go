package s3client

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type NewMinIOClientParams struct {
	Endpoint  string
	AccessKey string
	SecretKey string
}

func NewMinIOClient(p NewMinIOClientParams) *s3.Client {
	creds := credentials.NewStaticCredentialsProvider(p.AccessKey, p.SecretKey, "")
	
	sdkConfig, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion("us-east-1"),
		config.WithBaseEndpoint(p.Endpoint),
		config.WithCredentialsProvider(creds),
	)

	if err != nil {
		panic(err)
	}

	s3Client := s3.NewFromConfig(sdkConfig, func(o *s3.Options) {
		o.UsePathStyle = true
	})
	return s3Client
}