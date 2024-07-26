package get_appointment_by_id_uc

import (
	"context"
	"log/slog"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/entities"
	appointment_repository_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/repositories/appointment"
	get_appointment_by_id_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/appointment/get_appointment_by_id"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/middlewares/role"
)

type useCase struct {
	repository appointment_repository_contract.Repository
}

func NewUseCase(repository appointment_repository_contract.Repository) get_appointment_by_id_contract.UseCase {
	return &useCase{
		repository: repository,
	}
}

func (uc *useCase) Execute(ctx context.Context, userID, appointmentID uint, roleName role.Role) (*entities.Appointment, error) {
	slog.InfoContext(ctx, "getting appointment by id", "userId", userID, "appointmentId", appointmentID, "role", roleName)
	return uc.repository.GetByID(ctx, userID, appointmentID, roleName)
}
