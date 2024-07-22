package feedback_repository

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/entities"
	feedback_repository_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/repositories/feedback"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/infrastructure/shared/app_error"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/middlewares/role"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/persistence"
	"gorm.io/gorm"
)

type repository struct {
	dbService *persistence.DbService
}

func NewRepository(dbService *persistence.DbService) feedback_repository_contract.Repository {
	return &repository{
		dbService: dbService,
	}
}

func (rp *repository) GetByID(ctx context.Context, feedbackID uint) (*entities.Feedback, error) {
	tx := rp.dbService.Instance.WithContext(ctx)

	feedback := new(entities.Feedback)

	if err := tx.Where("id = ?", feedbackID).First(feedback).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, app_error.New(http.StatusNotFound, fmt.Sprintf("feedback with id %d not found", feedbackID))
		}

		return nil, err
	}

	return feedback, nil
}

func (rp *repository) GetByAppointmentID(ctx context.Context, appointmentID uint) (*entities.Feedback, error) {
	tx := rp.dbService.Instance.WithContext(ctx)

	feedback := new(entities.Feedback)

	if err := tx.Where("appointment_id = ?", appointmentID).First(feedback).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, app_error.New(http.StatusNotFound, fmt.Sprintf("feedback with appointment id %d not found", appointmentID))
		}

		return nil, err
	}

	return feedback, nil
}

func (rp *repository) List(ctx context.Context, userID, appointmentID uint, roleName role.Role) ([]entities.Feedback, error) {
	tx := rp.dbService.Instance.WithContext(ctx)

	feedbacks := new([]entities.Feedback)

	query := "appointments.patient_id = ? AND appointments.id = ?"

	if roleName == role.Doctor {
		query = "appointments.doctor_id = ? AND appointments.id = ?"
	}

	if err := tx.Preload("Appointment").
		Order("feedbacks.created_at DESC").
		Joins("JOIN appointments ON appointments.id = feedbacks.appointment_id").
		Where(query, userID, appointmentID).
		Find(&feedbacks).Error; err != nil {
		return nil, err
	}

	return *feedbacks, nil
}

func (rp *repository) Create(ctx context.Context, feedback *entities.Feedback) (*entities.Feedback, error) {
	tx := rp.dbService.Instance.WithContext(ctx)

	if err := tx.Create(feedback).Error; err != nil {
		return nil, err
	}

	return feedback, nil
}

func (rp *repository) Update(ctx context.Context, feedback *entities.Feedback) (*entities.Feedback, error) {
	tx := rp.dbService.Instance.WithContext(ctx)

	if err := tx.Model(feedback).Where("id = ?", feedback.ID).Updates(feedback).Error; err != nil {
		return nil, err
	}

	return feedback, nil
}

func (rp *repository) Delete(ctx context.Context, feedbackID uint) error {
	tx := rp.dbService.Instance.WithContext(ctx)

	if err := tx.Delete(&entities.Appointment{}, "id = ?", feedbackID).Error; err != nil {
		return err
	}

	return nil
}
