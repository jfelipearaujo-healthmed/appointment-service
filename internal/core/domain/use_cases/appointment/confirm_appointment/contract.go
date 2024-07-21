package confirm_appointment_contract

import (
	"context"
)

type UseCase interface {
	Execute(ctx context.Context, userID uint, appointmentID uint, confirmed bool) error
}
