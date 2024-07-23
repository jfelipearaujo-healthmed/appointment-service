package file_access_repository

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/entities"
	file_access_repository_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/repositories/file_access"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/infrastructure/shared/app_error"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/persistence"
	"gorm.io/gorm"
)

type repository struct {
	dbService *persistence.DbService
}

func NewRepository(dbService *persistence.DbService) file_access_repository_contract.Repository {
	return &repository{
		dbService: dbService,
	}
}

func (rp *repository) GetByID(ctx context.Context, userID uint, fileAccessID uint) (*entities.FileAccess, error) {
	tx := rp.dbService.Instance.WithContext(ctx)

	fileAccess := new(entities.FileAccess)

	if err := tx.Where("id = ? AND user_id = ?", fileAccessID, userID).First(fileAccess).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, app_error.New(http.StatusNotFound, fmt.Sprintf("file access with id %d not found", fileAccessID))
		}

		return nil, err
	}

	return fileAccess, nil
}

func (rp *repository) GetByIDForDoctor(ctx context.Context, userID, doctorID, fileAccessID uint) (*entities.FileAccess, error) {
	tx := rp.dbService.Instance.WithContext(ctx)

	fileAccess := new(entities.FileAccess)

	if err := tx.Where("id = ? AND user_id = ? AND doctor_id = ?", fileAccessID, userID, doctorID).First(fileAccess).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, app_error.New(http.StatusNotFound, fmt.Sprintf("file access with id %d not found", fileAccessID))
		}

		return nil, err
	}

	return fileAccess, nil
}

func (rp *repository) GetByAppointmentID(ctx context.Context, appointmentID uint) ([]entities.FileAccess, error) {
	tx := rp.dbService.Instance.WithContext(ctx)

	fileAccesses := new([]entities.FileAccess)

	if err := tx.Where("appointment_id = ? AND expired_at > NOW()", appointmentID).Find(&fileAccesses).Error; err != nil {
		return nil, err
	}

	return *fileAccesses, nil
}

func (rp *repository) List(ctx context.Context, userID uint) ([]entities.FileAccess, error) {
	tx := rp.dbService.Instance.WithContext(ctx)

	fileAccesses := new([]entities.FileAccess)

	if err := tx.Where("user_id = ?", userID).Find(&fileAccesses).Error; err != nil {
		return nil, err
	}

	return *fileAccesses, nil
}

func (rp *repository) Create(ctx context.Context, fileAccess *entities.FileAccess) (*entities.FileAccess, error) {
	tx := rp.dbService.Instance.WithContext(ctx)

	if err := tx.Create(fileAccess).Error; err != nil {
		return nil, err
	}

	return fileAccess, nil
}

func (rp *repository) Update(ctx context.Context, userID uint, fileAccess *entities.FileAccess) (*entities.FileAccess, error) {
	tx := rp.dbService.Instance.WithContext(ctx)

	if err := tx.Model(fileAccess).Where("id = ? AND user_id = ?", fileAccess.ID, userID).Updates(fileAccess).Error; err != nil {
		return nil, err
	}

	return fileAccess, nil
}

func (rp *repository) Delete(ctx context.Context, userID uint, fileAccessID uint) error {
	tx := rp.dbService.Instance.WithContext(ctx)

	if err := tx.Delete(&entities.FileAccess{}, "id = ? AND user_id = ?", fileAccessID, userID).Error; err != nil {
		return err
	}

	return nil
}
