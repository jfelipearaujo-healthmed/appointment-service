package storage

import (
	"context"
	"fmt"
	"log/slog"
	"mime/multipart"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewService(bucketName string, config aws.Config) StorageService {
	client := s3.NewFromConfig(config)

	return &Service{
		Client:        client,
		PresignClient: s3.NewPresignClient(client),
		BucketName:    bucketName,
	}
}

func (svc *Service) Upload(ctx context.Context, fileKey string, fileData multipart.File) (string, error) {
	output, err := svc.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(svc.BucketName),
		Key:    aws.String(fileKey),
		Body:   fileData,
	})
	if err != nil {
		return "", err
	}

	slog.InfoContext(ctx, "file uploaded", "file_key", fileKey, "location", output.ETag)

	return fmt.Sprintf("%s/%s", svc.BucketName, fileKey), nil
}

func (svc *Service) GetUrl(ctx context.Context, fileKey string) (string, error) {
	objectData, err := svc.PresignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(svc.BucketName),
		Key:    aws.String(fileKey),
	}, func(options *s3.PresignOptions) {
		options.Expires = time.Duration(24 * time.Hour)
	})
	if err != nil {
		return "", err
	}

	slog.InfoContext(ctx, "file url generated", "file_key", fileKey, "url", objectData.URL)

	return objectData.URL, nil
}
