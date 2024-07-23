package create_file_access_uc

import (
	"context"
	"net/http"
	"time"

	file_access_dto "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/dto/file_access"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/entities"
	appointment_repository_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/repositories/appointment"
	file_repository_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/repositories/file"
	file_access_repository_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/repositories/file_access"
	create_file_access_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/file_access/create_file_access"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/infrastructure/shared/app_error"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/middlewares/role"
)

const (
	dateTimeLayout = "2006-01-02 15:04"
)

type useCase struct {
	appointmentRepository appointment_repository_contract.Repository
	fileRepository        file_repository_contract.Repository
	fileAccessRepository  file_access_repository_contract.Repository
	location              *time.Location
}

func NewUseCase(
	appointmentRepository appointment_repository_contract.Repository,
	fileRepository file_repository_contract.Repository,
	fileAccessRepository file_access_repository_contract.Repository,
	location *time.Location,
) create_file_access_contract.UseCase {
	return &useCase{
		appointmentRepository: appointmentRepository,
		fileRepository:        fileRepository,
		fileAccessRepository:  fileAccessRepository,
		location:              location,
	}
}

func (uc *useCase) Execute(ctx context.Context, userID uint, request *file_access_dto.CreateFileAccess) (*entities.FileAccess, error) {
	_, err := uc.appointmentRepository.GetByID(ctx, userID, request.AppointmentID, role.Patient)
	if err != nil {
		return nil, err
	}

	_, err = uc.fileRepository.GetByID(ctx, userID, request.FileID)
	if err != nil {
		return nil, err
	}

	parsedExpiresAt, err := time.ParseInLocation(dateTimeLayout, request.ExpiresAt, uc.location)
	if err != nil {
		return nil, app_error.New(http.StatusBadRequest, "unable to parse the date and time provided")
	}

	year, month, day := parsedExpiresAt.Date()
	hour, minute, _ := parsedExpiresAt.Clock()

	expiresAt := time.Date(year, month, day, hour, minute, 0, 0, uc.location)

	if expiresAt.After(time.Now()) {
		return nil, app_error.New(http.StatusBadRequest, "date and time must be in the past")
	}

	fileAccess := &entities.FileAccess{
		UserID:        userID,
		FileID:        request.FileID,
		DoctorID:      request.DoctorID,
		AppointmentID: request.AppointmentID,
		ExpiresAt:     expiresAt,
	}

	existingFileAccess, err := uc.fileAccessRepository.GetByIDForDoctor(ctx, userID, request.DoctorID, request.FileID)
	if err != nil && !app_error.IsAppError(err) {
		return nil, err
	}

	if existingFileAccess != nil {
		return nil, app_error.New(http.StatusBadRequest, "file access already created")
	}

	fileAccess, err = uc.fileAccessRepository.Create(ctx, fileAccess)
	if err != nil {
		return nil, err
	}

	return fileAccess, nil
}
