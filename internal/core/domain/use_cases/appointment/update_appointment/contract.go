package update_appointment_contract

import (
	"context"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/dto/appointment_dto"
)

type UseCase interface {
	Execute(ctx context.Context, patientID, appointmentID uint, request *appointment_dto.UpdateAppointmentRequest) error
}
