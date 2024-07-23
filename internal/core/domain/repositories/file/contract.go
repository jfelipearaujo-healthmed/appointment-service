package file_repository_contract

import (
	"context"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/entities"
)

type Repository interface {
	GetByID(ctx context.Context, userID, fileID uint) (*entities.File, error)
	List(ctx context.Context, userID uint) ([]entities.File, error)
	Create(ctx context.Context, file *entities.File) (*entities.File, error)
	Update(ctx context.Context, userID uint, file *entities.File) (*entities.File, error)
	Delete(ctx context.Context, userID uint, fileID uint) error
}
