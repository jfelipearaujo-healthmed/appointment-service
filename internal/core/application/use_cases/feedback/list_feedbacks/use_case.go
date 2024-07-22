package list_feedbacks_uc

import (
	"context"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/entities"
	feedback_repository_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/repositories/feedback"
	list_feedbacks_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/feedback/list_feedbacks"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/middlewares/role"
)

type useCase struct {
	repository feedback_repository_contract.Repository
}

func NewUseCase(repository feedback_repository_contract.Repository) list_feedbacks_contract.UseCase {
	return &useCase{
		repository: repository,
	}
}

func (uc *useCase) Execute(ctx context.Context, userID, appointmentID uint, roleName role.Role) ([]entities.Feedback, error) {
	return uc.repository.List(ctx, userID, appointmentID, roleName)
}
