package list_appointments_contract

import (
	"context"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/entities"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/middlewares/role"
)

type UseCase interface {
	Execute(ctx context.Context, userID uint, roleName role.Role) ([]entities.Appointment, error)
}
