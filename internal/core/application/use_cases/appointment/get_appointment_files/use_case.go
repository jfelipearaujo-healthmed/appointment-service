package get_appointment_files_uc

import (
	"context"
	"fmt"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/entities"
	appointment_repository_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/repositories/appointment"
	file_repository_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/repositories/file"
	file_access_repository_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/repositories/file_access"
	get_appointment_files_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/appointment/get_appointment_files"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/middlewares/role"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/storage"
)

type useCase struct {
	appointmentRepository appointment_repository_contract.Repository
	fileAccessRepository  file_access_repository_contract.Repository
	fileRepository        file_repository_contract.Repository
	storageService        storage.StorageService
}

func NewUseCase(
	appointmentRepository appointment_repository_contract.Repository,
	fileAccessRepository file_access_repository_contract.Repository,
	fileRepository file_repository_contract.Repository,
	storageService storage.StorageService,
) get_appointment_files_contract.UseCase {
	return &useCase{
		appointmentRepository: appointmentRepository,
		fileAccessRepository:  fileAccessRepository,
		fileRepository:        fileRepository,
		storageService:        storageService,
	}
}

func (uc *useCase) Execute(ctx context.Context, userID, appointmentID uint) ([]entities.File, error) {
	_, err := uc.appointmentRepository.GetByID(ctx, userID, appointmentID, role.Doctor)
	if err != nil {
		return nil, err
	}

	filesAccess, err := uc.fileAccessRepository.GetByAppointmentID(ctx, appointmentID)
	if err != nil {
		return nil, err
	}

	fileIDs := make([]uint, 0)

	for _, fileAccess := range filesAccess {
		fileIDs = append(fileIDs, fileAccess.FileID)
	}

	files, err := uc.fileRepository.GetByFileIDs(ctx, fileIDs)
	if err != nil {
		return nil, err
	}

	for i, file := range files {
		fileKey := fmt.Sprintf("files/%d/%s", file.UserId, file.Name)

		url, err := uc.storageService.GetUrl(ctx, fileKey)
		if err != nil {
			return nil, err
		}

		files[i].Url = url
	}

	return files, nil
}
