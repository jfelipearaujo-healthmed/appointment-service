package appointment_repository_contract

import (
	"context"
	"time"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/entities"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/middlewares/role"
)

type Repository interface {
	GetByID(ctx context.Context, userID, appointmentID uint, roleName role.Role) (*entities.Appointment, error)
	GetByIDsAndDateTime(ctx context.Context, scheduleID uint, patientID, doctorID uint, dateTime time.Time) (*entities.Appointment, error)
	List(ctx context.Context, userID uint, roleName role.Role) ([]entities.Appointment, error)
	Create(ctx context.Context, appointment *entities.Appointment) (*entities.Appointment, error)
	Update(ctx context.Context, userID uint, appointment *entities.Appointment) (*entities.Appointment, error)
	Delete(ctx context.Context, userID uint, appointmentID uint) error
}
