package event_repository

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/entities"
	event_repository_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/repositories/event"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/infrastructure/shared/app_error"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/persistence"
	"gorm.io/gorm"
)

type repository struct {
	dbService *persistence.DbService
}

func NewRepository(dbService *persistence.DbService) event_repository_contract.Repository {
	return &repository{
		dbService: dbService,
	}
}

func (rp *repository) Create(ctx context.Context, event *entities.Event) (*entities.Event, error) {
	tx := rp.dbService.Instance.WithContext(ctx)

	if err := tx.Create(event).Error; err != nil {
		return nil, err
	}

	return event, nil
}

func (rp *repository) GetByIDsAndDateTime(ctx context.Context, scheduleID uint, patientID uint, doctorID uint, dateTime time.Time) (*entities.Event, error) {
	tx := rp.dbService.Instance.WithContext(ctx)

	event := new(entities.Event)

	query := tx.Where("schedule_id = ?", scheduleID)
	query = query.Where("patient_id = ?", patientID)
	query = query.Where("doctor_id = ?", doctorID)
	query = query.Where("date_time = ?", dateTime)
	query = query.Where("outcome IS NULL")

	result := query.First(event)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, app_error.New(http.StatusNotFound, fmt.Sprintf("event with schedule id %d, patient id %d, doctor id %d and date time %s not found", scheduleID, patientID, doctorID, dateTime))
		}

		return nil, result.Error
	}

	return event, nil
}
