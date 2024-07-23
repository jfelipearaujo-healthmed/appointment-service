package storage

import (
	"context"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Service struct {
	Client        *s3.Client
	PresignClient *s3.PresignClient
	BucketName    string
}

type StorageService interface {
	Upload(ctx context.Context, fileKey string, fileData multipart.File) (string, error)
	GetUrl(ctx context.Context, fileKey string) (string, error)
}
