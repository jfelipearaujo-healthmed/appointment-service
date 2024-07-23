package get_file_by_id_contract

import (
	"context"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/entities"
)

type UseCase interface {
	Execute(ctx context.Context, userID, fileID uint) (*entities.File, error)
}
