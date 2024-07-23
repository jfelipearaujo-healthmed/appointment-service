package list_files_uc

import (
	"context"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/entities"
	file_repository_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/repositories/file"
	list_files_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/file/list_files"
)

type useCase struct {
	repository file_repository_contract.Repository
}

func New(repository file_repository_contract.Repository) list_files_contract.UseCase {
	return &useCase{
		repository: repository,
	}
}

func (u *useCase) Execute(ctx context.Context, userID uint) ([]entities.File, error) {
	return u.repository.List(ctx, userID)
}
