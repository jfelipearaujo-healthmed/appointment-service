package list_file_access_contract

import (
	"context"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/entities"
)

type UseCase interface {
	Execute(ctx context.Context, userID uint) ([]entities.FileAccess, error)
}
