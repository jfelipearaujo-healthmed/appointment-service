package get_feedback_by_id_contract

import (
	"context"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/entities"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/middlewares/role"
)

type UseCase interface {
	Execute(ctx context.Context, userID, appointmentID, feedbackID uint, roleName role.Role) (*entities.Feedback, error)
}
