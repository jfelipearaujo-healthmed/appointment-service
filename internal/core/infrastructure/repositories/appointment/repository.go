package appointment_repository

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/entities"
	appointment_repository_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/repositories/appointment"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/infrastructure/shared/app_error"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/cache"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/middlewares/role"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/persistence"
	"gorm.io/gorm"
)

const (
	cacheKeyPrefix string        = "appointment:%d"
	ttl            time.Duration = time.Hour * 24
)

type repository struct {
	cache     cache.Cache
	dbService *persistence.DbService
}

func NewRepository(cache cache.Cache, dbService *persistence.DbService) appointment_repository_contract.Repository {
	return &repository{
		cache:     cache,
		dbService: dbService,
	}
}

func (rp *repository) GetByID(ctx context.Context, userID uint, appointmentID uint, roleName role.Role) (*entities.Appointment, error) {
	cacheKey := fmt.Sprintf(cacheKeyPrefix, appointmentID)
	return cache.WithCache(ctx, rp.cache, cacheKey, ttl, func() (*entities.Appointment, error) {
		tx := rp.dbService.Instance.WithContext(ctx)

		appointment := new(entities.Appointment)

		query := "patient_id = ? AND id = ?"
		if roleName == role.Doctor {
			query = "doctor_id = ? AND id = ?"
		}

		result := tx.Where(query, userID, appointmentID).First(appointment)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return nil, app_error.New(http.StatusNotFound, fmt.Sprintf("appointment with id %d not found", appointmentID))
			}

			return nil, result.Error
		}

		return appointment, nil
	})
}

func (rp *repository) GetByIDsAndDateTime(ctx context.Context, scheduleID uint, patientID uint, doctorID uint, dateTime time.Time) (*entities.Appointment, error) {
	tx := rp.dbService.Instance.WithContext(ctx)

	appointment := new(entities.Appointment)
	result := tx.Where("schedule_id = ? AND patient_id = ? AND doctor_id = ? AND date_time = ?", scheduleID, patientID, doctorID, dateTime).First(appointment)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, app_error.New(http.StatusNotFound, fmt.Sprintf("appointment with schedule id %d, patient id %d, doctor id %d and date time %s not found", scheduleID, patientID, doctorID, dateTime))
		}

		return nil, result.Error
	}

	return appointment, nil
}

func (rp *repository) List(ctx context.Context, userID uint, roleName role.Role) ([]entities.Appointment, error) {
	tx := rp.dbService.Instance.WithContext(ctx)

	appointments := new([]entities.Appointment)

	query := "patient_id = ?"

	if roleName == role.Doctor {
		query = "doctor_id = ?"
	}

	if err := tx.Where(query, userID).Find(&appointments).Error; err != nil {
		return nil, err
	}

	return *appointments, nil
}

func (rp *repository) Create(ctx context.Context, appointment *entities.Appointment) (*entities.Appointment, error) {
	tx := rp.dbService.Instance.WithContext(ctx)

	if err := tx.Create(appointment).Error; err != nil {
		return nil, err
	}

	return appointment, nil
}

func (rp *repository) Update(ctx context.Context, userID uint, appointment *entities.Appointment) (*entities.Appointment, error) {
	tx := rp.dbService.Instance.WithContext(ctx)

	if err := tx.Model(appointment).Where("patient_id = ? AND id = ?", userID, appointment.ID).Updates(appointment).Error; err != nil {
		return nil, err
	}

	return appointment, nil
}

func (rp *repository) Delete(ctx context.Context, userID uint, appointmentID uint) error {
	tx := rp.dbService.Instance.WithContext(ctx)

	if err := tx.Delete(&entities.Appointment{}, "patient_id = ? AND id = ?", userID, appointmentID).Error; err != nil {
		return err
	}

	return nil
}
