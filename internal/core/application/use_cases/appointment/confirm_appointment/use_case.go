package confirm_appointment_uc

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/entities"
	appointment_repository_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/repositories/appointment"
	confirm_appointment_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/appointment/confirm_appointment"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/infrastructure/shared/app_error"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/middlewares/role"
)

type useCase struct {
	repository appointment_repository_contract.Repository
}

func NewUseCase(repository appointment_repository_contract.Repository) confirm_appointment_contract.UseCase {
	return &useCase{
		repository: repository,
	}
}

func (uc *useCase) Execute(ctx context.Context, userID uint, appointmentID uint, confirmed bool) error {
	slog.InfoContext(ctx, "confirming appointment", "userId", userID, "appointmentId", appointmentID, "confirmed", confirmed)

	appointment, err := uc.repository.GetByID(ctx, userID, appointmentID, role.Doctor)
	if err != nil {
		return err
	}

	slog.InfoContext(ctx, "appointment retrieved", "appointment", appointment)

	if appointment.Status != entities.ScheduleInAnalysis &&
		appointment.Status != entities.ReScheduleInAnalysis &&
		appointment.Status != entities.WaitingForConfirmation {
		return app_error.New(http.StatusBadRequest, "appointment is not in schedule or re-schedule status")
	}

	if appointment.Status == entities.Confirmed {
		return app_error.New(http.StatusBadRequest, "appointment is already confirmed")
	}

	now := time.Now()

	if confirmed {
		slog.InfoContext(ctx, "confirming appointment", "userId", userID, "appointmentId", appointmentID, "confirmed", confirmed)
		appointment.Status = entities.Confirmed
		appointment.ConfirmedAt = &now
	} else {
		slog.InfoContext(ctx, "declining appointment", "userId", userID, "appointmentId", appointmentID, "confirmed", confirmed)
		reason := "Doctor does not confirmed this appointment, please reschedule"
		appointment.Status = entities.Cancelled
		appointment.CancelledAt = &now
		appointment.CancelledBy = &userID
		appointment.CancelledReason = &reason
	}

	slog.InfoContext(ctx, "updating appointment", "appointment", appointment)

	if _, err := uc.repository.Update(ctx, userID, appointment); err != nil {
		return err
	}

	slog.InfoContext(ctx, "appointment updated", "appointment", appointment)

	return nil
}
