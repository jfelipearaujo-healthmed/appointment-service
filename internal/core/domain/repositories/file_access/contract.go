package file_access_repository_contract

import (
	"context"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/entities"
)

type Repository interface {
	GetByID(ctx context.Context, userID, fileAccessID uint) (*entities.FileAccess, error)
	GetByIDForDoctor(ctx context.Context, userID, doctorID, fileAccessID uint) (*entities.FileAccess, error)
	List(ctx context.Context, userID uint) ([]entities.FileAccess, error)
	Create(ctx context.Context, fileAccess *entities.FileAccess) (*entities.FileAccess, error)
	Update(ctx context.Context, userID uint, fileAccess *entities.FileAccess) (*entities.FileAccess, error)
	Delete(ctx context.Context, userID uint, fileAccessID uint) error
}
