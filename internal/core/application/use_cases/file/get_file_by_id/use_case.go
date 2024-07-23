package get_file_by_id_uc

import (
	"context"
	"fmt"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/entities"
	file_repository_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/repositories/file"
	get_file_by_id_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/file/get_file_by_id"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/storage"
)

type useCase struct {
	repository     file_repository_contract.Repository
	storageService storage.StorageService
}

func NewUseCase(
	repository file_repository_contract.Repository,
	storageService storage.StorageService,
) get_file_by_id_contract.UseCase {
	return &useCase{
		repository:     repository,
		storageService: storageService,
	}
}

func (uc *useCase) Execute(ctx context.Context, userID, fileID uint) (*entities.File, error) {
	file, err := uc.repository.GetByID(ctx, userID, fileID)
	if err != nil {
		return nil, err
	}

	fileKey := fmt.Sprintf("files/%d/%s", userID, file.Name)

	url, err := uc.storageService.GetUrl(ctx, fileKey)
	if err != nil {
		return nil, err
	}

	file.Url = url

	return file, nil
}
