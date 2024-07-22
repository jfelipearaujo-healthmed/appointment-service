package get_feedback_by_id_uc

import (
	"context"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/entities"
	feedback_repository_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/repositories/feedback"
	get_feedback_by_id_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/feedback/get_feedback_by_id"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/middlewares/role"
)

type useCase struct {
	repository feedback_repository_contract.Repository
}

func NewUseCase(repository feedback_repository_contract.Repository) get_feedback_by_id_contract.UseCase {
	return &useCase{
		repository: repository,
	}
}

func (uc *useCase) Execute(ctx context.Context, userID, appointmentID, feedbackID uint, roleName role.Role) (*entities.Feedback, error) {
	return uc.repository.GetByID(ctx, userID, appointmentID, feedbackID, roleName)
}
