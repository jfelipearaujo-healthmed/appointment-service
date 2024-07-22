package feedback_repository_contract

import (
	"context"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/entities"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/middlewares/role"
)

type Repository interface {
	GetByID(ctx context.Context, feedbackID uint) (*entities.Feedback, error)
	GetByAppointmentID(ctx context.Context, appointmentID uint) (*entities.Feedback, error)
	List(ctx context.Context, userID uint, roleName role.Role) ([]entities.Feedback, error)
	Create(ctx context.Context, feedback *entities.Feedback) (*entities.Feedback, error)
	Update(ctx context.Context, feedback *entities.Feedback) (*entities.Feedback, error)
	Delete(ctx context.Context, feedbackID uint) error
}
