package upload_file_contract

import (
	"context"
	"mime/multipart"
)

type UseCase interface {
	Execute(ctx context.Context, userID uint, fileName, mimeType string, fileSize int64, fileData multipart.File) error
}
