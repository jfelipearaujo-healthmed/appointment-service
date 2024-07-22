package cancel_appointment_uc

import (
	"context"

	appointment_repository_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/repositories/appointment"
	cancel_appointment_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/appointment/cancel_appointment"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/middlewares/role"
)

type useCase struct {
	repository appointment_repository_contract.Repository
}

func NewUseCase(
	repository appointment_repository_contract.Repository,
) cancel_appointment_contract.UseCase {
	return &useCase{
		repository: repository,
	}
}

func (uc *useCase) Execute(ctx context.Context, userID uint, appointmentID uint, roleName role.Role, reason string) error {
	appointment, err := uc.repository.GetByID(ctx, userID, appointmentID, roleName)
	if err != nil {
		return err
	}

	appointment.Cancel(userID, reason)

	_, err = uc.repository.Update(ctx, userID, appointment)
	if err != nil {
		return err
	}

	return nil
}
