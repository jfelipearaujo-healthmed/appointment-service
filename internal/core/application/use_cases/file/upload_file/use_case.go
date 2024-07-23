package upload_file_uc

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/entities"
	file_repository_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/repositories/file"
	upload_file_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/file/upload_file"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/storage"
)

type useCase struct {
	storageService storage.StorageService
	fileRepository file_repository_contract.Repository
}

func NewUseCase(
	storageService storage.StorageService,
	fileRepository file_repository_contract.Repository,
) upload_file_contract.UseCase {
	return &useCase{
		storageService: storageService,
		fileRepository: fileRepository,
	}
}

func (uc *useCase) Execute(ctx context.Context, userID uint, fileName, mimeType string, fileSize int64, fileData multipart.File) error {
	randomName := fmt.Sprintf("%s-%s", uuid.NewString(), fileName)
	fileKey := fmt.Sprintf("files/%d/%s", userID, randomName)

	fileUrl, err := uc.storageService.Upload(ctx, fileKey, fileData)
	if err != nil {
		return err
	}

	file := &entities.File{
		UserId:       userID,
		Name:         randomName,
		OriginalName: fileName,
		MimeType:     mimeType,
		Size:         fileSize,
		Url:          fileUrl,
	}

	_, err = uc.fileRepository.Create(ctx, file)

	return err
}
