package cancel_appointment_uc

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/entities"
	appointment_repository_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/repositories/appointment"
	cancel_appointment_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/appointment/cancel_appointment"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/infrastructure/shared/app_error"
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
	slog.InfoContext(ctx, "canceling appointment", "userId", userID, "appointmentId", appointmentID, "role", roleName, "reason", reason)

	appointment, err := uc.repository.GetByID(ctx, userID, appointmentID, roleName)
	if err != nil {
		return err
	}

	slog.InfoContext(ctx, "appointment retrieved", "appointment", appointment)

	if appointment.Status == entities.Cancelled {
		return app_error.New(http.StatusBadRequest, "appointment is already cancelled")
	}

	appointment.Cancel(userID, reason)

	slog.InfoContext(ctx, "updating appointment", "appointment", appointment)

	if _, err := uc.repository.Update(ctx, userID, appointment); err != nil {
		return err
	}

	slog.InfoContext(ctx, "appointment updated", "appointment", appointment)

	return nil
}
