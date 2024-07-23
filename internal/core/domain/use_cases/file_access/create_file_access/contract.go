package create_file_access_contract

import (
	"context"

	file_access_dto "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/dto/file_access"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/entities"
)

type UseCase interface {
	Execute(ctx context.Context, userID uint, request *file_access_dto.CreateFileAccess) (*entities.FileAccess, error)
}
