package create_feedback_contract

import (
	"context"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/dto/feedback_dto"
)

type UseCase interface {
	Execute(ctx context.Context, patientID, appointmentID uint, request *feedback_dto.CreateFeedbackRequest) error
}
