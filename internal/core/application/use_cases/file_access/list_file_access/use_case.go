package list_file_access_uc

import (
	"context"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/entities"
	file_access_repository_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/repositories/file_access"
	list_file_access_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/file_access/list_file_access"
)

type useCase struct {
	repository file_access_repository_contract.Repository
}

func NewUseCase(repository file_access_repository_contract.Repository) list_file_access_contract.UseCase {
	return &useCase{
		repository: repository,
	}
}

func (uc *useCase) Execute(ctx context.Context, userID uint) ([]entities.FileAccess, error) {
	return uc.repository.List(ctx, userID)
}
