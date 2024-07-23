package file_repository

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/entities"
	file_repository_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/repositories/file"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/infrastructure/shared/app_error"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/persistence"
	"gorm.io/gorm"
)

type repository struct {
	dbService *persistence.DbService
}

func NewRepository(dbService *persistence.DbService) file_repository_contract.Repository {
	return &repository{
		dbService: dbService,
	}
}

func (rp *repository) GetByID(ctx context.Context, userID uint, fileID uint) (*entities.File, error) {
	tx := rp.dbService.Instance.WithContext(ctx)

	file := new(entities.File)

	if err := tx.Where("id = ? AND user_id = ?", fileID, userID).First(file).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, app_error.New(http.StatusNotFound, fmt.Sprintf("file with id %d not found", fileID))
		}

		return nil, err
	}

	return file, nil
}

func (rp *repository) GetByFileIDs(ctx context.Context, fileIDs []uint) ([]entities.File, error) {
	tx := rp.dbService.Instance.WithContext(ctx)

	files := new([]entities.File)

	if err := tx.Where("id IN (?)", fileIDs).Find(&files).Error; err != nil {
		return nil, err
	}

	return *files, nil
}

func (rp *repository) List(ctx context.Context, userID uint) ([]entities.File, error) {
	tx := rp.dbService.Instance.WithContext(ctx)

	files := new([]entities.File)

	if err := tx.Where("user_id = ?", userID).Find(&files).Error; err != nil {
		return nil, err
	}

	return *files, nil
}

func (rp *repository) Create(ctx context.Context, file *entities.File) (*entities.File, error) {
	tx := rp.dbService.Instance.WithContext(ctx)

	if err := tx.Create(file).Error; err != nil {
		return nil, err
	}

	return file, nil
}

func (rp *repository) Update(ctx context.Context, userID uint, file *entities.File) (*entities.File, error) {
	tx := rp.dbService.Instance.WithContext(ctx)

	if err := tx.Model(file).Where("id = ? AND user_id = ?", file.ID, userID).Updates(file).Error; err != nil {
		return nil, err
	}

	return file, nil
}

func (rp *repository) Delete(ctx context.Context, userID uint, fileID uint) error {
	tx := rp.dbService.Instance.WithContext(ctx)

	if err := tx.Delete(&entities.File{}, "id = ? AND user_id = ?", fileID, userID).Error; err != nil {
		return err
	}

	return nil
}
