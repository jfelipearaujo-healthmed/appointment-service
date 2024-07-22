package cancel_appointment_contract

import (
	"context"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/middlewares/role"
)

type UseCase interface {
	Execute(ctx context.Context, userID uint, appointmentID uint, roleName role.Role, reason string) error
}
