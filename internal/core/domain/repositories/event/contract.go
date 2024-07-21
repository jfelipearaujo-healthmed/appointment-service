package event_repository_contract

import (
	"context"
	"time"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/entities"
)

type Repository interface {
	Create(ctx context.Context, event *entities.Event) (*entities.Event, error)
	GetByIDsAndDateTime(ctx context.Context, scheduleID uint, patientID, doctorID uint, dateTime time.Time) (*entities.Event, error)
}
