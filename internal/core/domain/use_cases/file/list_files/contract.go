package list_files_contract

import (
	"context"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/entities"
)

type UseCase interface {
	Execute(ctx context.Context, userID uint) ([]entities.File, error)
}
