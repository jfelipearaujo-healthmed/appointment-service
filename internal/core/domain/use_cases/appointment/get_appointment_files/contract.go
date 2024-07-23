package get_appointment_files_contract

import (
	"context"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/entities"
)

type UseCase interface {
	Execute(ctx context.Context, userID, appointmentID uint) ([]entities.File, error)
}
