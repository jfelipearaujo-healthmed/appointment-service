package event_repository_contract

import (
	"context"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/entities"
)

type Repository interface {
	Create(ctx context.Context, event *entities.Event) (*entities.Event, error)
	GetByIDsAndDateTime(ctx context.Context, event *entities.Event) (*entities.Event, error)
}
